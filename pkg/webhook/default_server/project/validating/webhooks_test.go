/*
Copyright 2018 Pressinfra SRL.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validating

import (
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/presslabs/dashboard/pkg/internal/projectns"
)

var _ = Describe("Project webhook", func() {
	var (
		// stop channel for controller manager
		stop chan struct{}

		webhook *NamespaceCreateHandler

		proj *projectns.ProjectNamespace

		organizationName string
	)

	BeforeEach(func() {
		organizationName = fmt.Sprintf("organization%d", rand.Int31())
		projectName := fmt.Sprintf("project%d", rand.Int31())

		proj = projectns.New(
			&corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: projectName,
				},
			})
		webhook = &NamespaceCreateHandler{}

		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).NotTo(HaveOccurred())

		webhook.InjectClient(mgr.GetClient())
		webhook.InjectDecoder(mgr.GetAdmissionDecoder())

		stop = StartTestManager(mgr)
	})

	AfterEach(func() {
		close(stop)
	})

	It("returns error when metadata is missing", func() {
		proj.SetLabels(map[string]string{
			"presslabs.com/kind": "project",
		})

		allowed, reason, err := webhook.validatingNamespaceFn(proj)

		Expect(allowed).To(Equal(false))
		Expect(reason).To(Equal("validation failed"))
		Expect(err).To(MatchError(ContainSubstring("required label \"presslabs.com/organization\" is missing")))
		Expect(err).To(MatchError(ContainSubstring("required label \"presslabs.com/project\" is missing")))
		Expect(err).To(MatchError(ContainSubstring("required annotation \"presslabs.com/created-by\" is missing")))
	})
	It("doesn't validate when kind is not a project", func() {
		proj.SetLabels(map[string]string{
			"presslabs.com/kind": "not-a-project",
		})

		allowed, reason, err := webhook.validatingNamespaceFn(proj)

		Expect(allowed).To(Equal(true))
		Expect(reason).To(Equal("not a project, skipping validation"))
		Expect(err).To(BeNil())
	})
	It("doesn't return error if metadata is provided", func() {
		proj.SetLabels(map[string]string{
			"presslabs.com/organization": organizationName,
			"presslabs.com/project":      proj.Name,
			"presslabs.com/kind":         "project",
		})

		proj.SetAnnotations(map[string]string{
			"presslabs.com/created-by": "Andi",
		})

		allowed, reason, err := webhook.validatingNamespaceFn(proj)
		Expect(allowed).To(Equal(true))
		Expect(reason).To(Equal("allowed to be admitted"))
		Expect(err).To(BeNil())
	})
})
