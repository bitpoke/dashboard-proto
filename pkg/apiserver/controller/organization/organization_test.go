/*
Copyright 2018 Pressinfra SRL

This file is subject to the terms and conditions defined in file LICENSE,
which is part of this source code package.
*/

package organization

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gosimple/slug"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/presslabs/dashboard/pkg/controller"
	. "github.com/presslabs/dashboard/pkg/internal/testutil/gomega"

	orgv1 "github.com/presslabs/dashboard-go/pkg/proto/presslabs/dashboard/organizations/v1"
	"github.com/presslabs/dashboard/pkg/apiserver/internal/metadata"
	"github.com/presslabs/dashboard/pkg/internal/organization"
)

const (
	ctxTimeout    = time.Second * 3
	deleteTimeout = time.Second
	updateTimeout = time.Second
)

// createOrganization is a helper func that creates an organization
func createOrganization(name, displayName, userID string) *organization.Organization {
	org := organization.New(&corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: organization.NamespaceName(name),
			Labels: map[string]string{
				"presslabs.com/kind":         "organization",
				"presslabs.com/organization": name,
			},
			Annotations: map[string]string{
				"presslabs.com/created-by": userID,
			},
		},
	})
	org.UpdateDisplayName(displayName)

	return org
}

// getNamespaceFn is a helper func that returns an organization
func getNamespaceFn(ctx context.Context, c client.Client, key client.ObjectKey) func() corev1.Namespace {
	return func() corev1.Namespace {
		var orgNs corev1.Namespace
		Expect(c.Get(ctx, key, &orgNs)).To(Succeed())
		return orgNs
	}
}

func expectProperNamespace(c client.Client, name, displayName, userID string) {
	var ns corev1.Namespace
	key := client.ObjectKey{
		Name: organization.NamespaceName(name),
	}
	Expect(c.Get(context.TODO(), key, &ns)).To(Succeed())
	Expect(ns.Name).To(Equal(fmt.Sprintf("org-%s", name)))
	Expect(ns.Labels).To(HaveKeyWithValue("presslabs.com/kind", "organization"))
	Expect(ns.Labels).To(HaveKeyWithValue("presslabs.com/organization", name))
	Expect(ns.Annotations).To(HaveKeyWithValue("presslabs.com/display-name", displayName))
	Expect(ns.Annotations).To(HaveKeyWithValue("presslabs.com/created-by", userID))
}

var _ = Describe("API server", func() {
	var (
		// stop channel for apiserver
		stop chan struct{}
		// controller k8s client
		c client.Client
		// client connection to an RPC server
		conn *grpc.ClientConn
		// orgClient
		orgClient orgv1.OrganizationsServiceClient
	)

	var (
		id, autoID     string
		name, autoName string
		displayName    string
		userID         string
	)

	BeforeEach(func() {
		mgr, err := manager.New(cfg, manager.Options{})
		Expect(err).To(Succeed())

		server := SetupAPIServer(mgr)
		// add ourselves to the server
		Add(server)

		// create new k8s client
		c, err = client.New(cfg, client.Options{})
		Expect(err).To(Succeed())

		// Add controllers for testing side effects
		Expect(controller.AddToManager(mgr)).To(Succeed())

		stop = StartTestManager(mgr)

		conn, err = grpc.Dial(server.GetGRPCAddr(), grpc.WithInsecure(), grpc.WithBlock(),
			grpc.WithTimeout(ctxTimeout))
		Expect(err).To(Succeed())

		orgClient = orgv1.NewOrganizationsServiceClient(conn)

		name = fmt.Sprintf("%d", rand.Int31())
		id = fmt.Sprintf("orgs/%s", name)
		displayName = fmt.Sprintf("Org %s Inc", name)
		autoName = slug.Make(displayName)
		autoID = fmt.Sprintf("orgs/%s", autoName)
		userID = fmt.Sprintf("user#%s", name)
		metadata.FakeSubject = userID
	})

	AfterEach(func() {
		// close the gRPC client connection
		conn.Close()
		// stop the manager and API server
		close(stop)

		// delete k8s resources
		orgs := &corev1.NamespaceList{}
		opts := &client.ListOptions{}
		err := opts.SetLabelSelector(fmt.Sprintf("presslabs.com/kind=organization"))
		Expect(err).To(Succeed())

		err = c.List(context.TODO(), opts, orgs)
		Expect(err).To(Succeed())

		for _, org := range orgs.Items {
			if org.Status.Phase == corev1.NamespaceTerminating {
				continue
			}

			err = c.Delete(context.TODO(), &org)
			Expect(err).To(Succeed())

			key, err := client.ObjectKeyFromObject(&org)
			Expect(err).To(Succeed())
			Eventually(getNamespaceFn(context.TODO(), c, key), deleteTimeout).Should(
				BeInPhase(corev1.NamespaceTerminating))
		}
	})

	Describe("at Create request", func() {
		It("returns AlreadyExists error when organization already exists", func() {
			org := createOrganization(name, displayName, userID)
			Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			req := orgv1.CreateOrganizationRequest{
				Organization: orgv1.Organization{
					Name:        id,
					DisplayName: displayName,
				},
			}

			_, err := orgClient.CreateOrganization(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.AlreadyExists))
		})

		It("returns error when no name is given", func() {
			req := orgv1.CreateOrganizationRequest{
				Organization: orgv1.Organization{},
			}
			_, err := orgClient.CreateOrganization(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when name is not fully qualified", func() {
			req := orgv1.CreateOrganizationRequest{
				Organization: orgv1.Organization{
					Name: "not-fully-qualified",
				},
			}
			_, err := orgClient.CreateOrganization(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("returns error when name is empty", func() {
			req := orgv1.CreateOrganizationRequest{
				Organization: orgv1.Organization{
					Name: "orgs/",
				},
			}
			_, err := orgClient.CreateOrganization(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.InvalidArgument))
		})

		It("creates organization when no organization name is given", func() {
			req := orgv1.CreateOrganizationRequest{
				Organization: orgv1.Organization{
					DisplayName: displayName,
				},
			}
			resp, err := orgClient.CreateOrganization(context.TODO(), &req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Name).To(Equal(autoID))
			expectProperNamespace(c, slug.Make(displayName), displayName, userID)
		})

		It("creates organization when name is given", func() {
			req := orgv1.CreateOrganizationRequest{
				Organization: orgv1.Organization{
					Name:        id,
					DisplayName: displayName,
				},
			}
			resp, err := orgClient.CreateOrganization(context.TODO(), &req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Name).To(Equal(id))
			expectProperNamespace(c, name, displayName, userID)
		})

		It("fills display_name when no one is given", func() {
			req := orgv1.CreateOrganizationRequest{
				Organization: orgv1.Organization{
					Name: id,
				},
			}
			resp, err := orgClient.CreateOrganization(context.TODO(), &req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Name).To(Equal(id))
			expectProperNamespace(c, name, name, userID)
		})
	})

	Describe("at Get request", func() {
		It("returns the organization", func() {
			org := createOrganization(name, displayName, userID)
			Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			req := orgv1.GetOrganizationRequest{
				Name: id,
			}

			resp, err := orgClient.GetOrganization(context.TODO(), &req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Name).To(Equal(id))
			Expect(resp.DisplayName).To(Equal(displayName))
		})

		It("returns NotFound when organization does not exist", func() {
			req := orgv1.GetOrganizationRequest{
				Name: id,
			}
			_, err := orgClient.GetOrganization(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})

		It("returns NotFound when organization namespace is not active", func() {
			termOrgName := fmt.Sprintf("%s-terminating", name)
			termOrg := createOrganization(termOrgName, displayName, userID)
			termOrg.Unwrap().ObjectMeta.Finalizers = []string{"api.dashboard.presslabs.org/terminating-org-namespace"}
			Expect(c.Create(context.TODO(), termOrg.Unwrap())).To(Succeed())

			key := types.NamespacedName{Name: organization.NamespaceName(termOrgName)}
			Expect(c.Delete(context.TODO(), termOrg.Unwrap())).To(Succeed())
			Eventually(getNamespaceFn(context.TODO(), c, key), deleteTimeout).Should(
				BeInPhase(corev1.NamespaceTerminating))

			req := orgv1.GetOrganizationRequest{
				Name: fmt.Sprintf("orgs/%s", termOrgName),
			}
			_, err := orgClient.GetOrganization(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})
	})

	Describe("at Delete request", func() {
		It("deletes existing organization", func() {
			org := createOrganization(name, displayName, userID)
			Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			req := orgv1.DeleteOrganizationRequest{
				Name: id,
			}
			_, err := orgClient.DeleteOrganization(context.TODO(), &req)
			Expect(err).ToNot(HaveOccurred())

			key := client.ObjectKey{
				Name: organization.NamespaceName(name),
			}

			Eventually(getNamespaceFn(context.TODO(), c, key), deleteTimeout).Should(
				BeInPhase(corev1.NamespaceTerminating))
		})

		It("returns NotFound when organization does not exist", func() {
			req := orgv1.DeleteOrganizationRequest{
				Name: id,
			}
			_, err := orgClient.DeleteOrganization(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})
	})

	Describe("at Update request", func() {
		BeforeEach(func() {
			org := createOrganization(name, displayName, userID)
			Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
		})
		It("updates display_name of existing organization", func() {
			newDisplayName := "The New Display Name"
			req := orgv1.UpdateOrganizationRequest{
				Organization: orgv1.Organization{
					Name:        id,
					DisplayName: newDisplayName,
				},
			}
			_, err := orgClient.UpdateOrganization(context.TODO(), &req)
			Expect(err).ToNot(HaveOccurred())

			key := client.ObjectKey{
				Name: organization.NamespaceName(name),
			}

			Eventually(getNamespaceFn(context.TODO(), c, key), updateTimeout).Should(
				HaveAnnotation("presslabs.com/display-name", newDisplayName))
		})
		It("sets display_name to default when no one is given", func() {
			newDisplayName := ""
			req := orgv1.UpdateOrganizationRequest{
				Organization: orgv1.Organization{
					Name:        id,
					DisplayName: newDisplayName,
				},
			}
			_, err := orgClient.UpdateOrganization(context.TODO(), &req)
			Expect(err).ToNot(HaveOccurred())

			key := client.ObjectKey{
				Name: organization.NamespaceName(name),
			}

			Eventually(getNamespaceFn(context.TODO(), c, key), updateTimeout).Should(
				HaveAnnotation("presslabs.com/display-name", name))
		})

		It("returns NotFound when organization does not exist", func() {
			req := orgv1.UpdateOrganizationRequest{
				Organization: orgv1.Organization{
					Name: "orgs/inexistent",
				},
			}
			_, err := orgClient.UpdateOrganization(context.TODO(), &req)
			Expect(status.Convert(err).Code()).To(Equal(codes.NotFound))
		})
	})

	Describe("at List request", func() {
		var myOrgsCount = 3
		BeforeEach(func() {
			for i := 1; i <= myOrgsCount; i++ {
				_name := fmt.Sprintf("%s-%02d", name, i)
				_displayName := fmt.Sprintf("%s %02d Inc.", name, i)
				org := createOrganization(_name, _displayName, userID)
				Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
			}
			org := createOrganization(name, displayName, "user#another")
			Expect(c.Create(context.TODO(), org.Unwrap())).To(Succeed())
		})

		It("returns only my organizations", func() {
			req := orgv1.ListOrganizationsRequest{}
			Eventually(func() ([]orgv1.Organization, error) {
				resp, err := orgClient.ListOrganizations(context.TODO(), &req)
				return resp.Organizations, err
			}).Should(HaveLen(myOrgsCount))
		})

		It("returns only active organizations", func() {
			termOrg := createOrganization("terminating-org", displayName, userID)
			termOrg.Status.Phase = corev1.NamespaceTerminating
			Expect(c.Create(context.TODO(), termOrg.Unwrap())).To(Succeed())

			req := orgv1.ListOrganizationsRequest{}
			Eventually(func() ([]orgv1.Organization, error) {
				resp, err := orgClient.ListOrganizations(context.TODO(), &req)
				return resp.Organizations, err
			}).Should(HaveLen(myOrgsCount))
		})
	})
})