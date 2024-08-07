package main

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types" // Correct import for NamespacedName
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controller Suite")
}

var _ = Describe("SecurityPolicyAuditor Controller", func() {
	var (
		r      *reconciler
		client client.Client
		scheme *runtime.Scheme
	)

	BeforeEach(func() {
		scheme = runtime.NewScheme()
		Expect(corev1.AddToScheme(scheme)).NotTo(HaveOccurred())
		Expect(rbacv1.AddToScheme(scheme)).NotTo(HaveOccurred())

		client = fake.NewClientBuilder().WithScheme(scheme).Build()
		r = &reconciler{
			client: client,
			scheme: scheme,
		}
	})

	Context("When creating a ConfigMap", func() {
		It("should audit the ConfigMap for sensitive data", func() {
			ctx := context.Background()

			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
				Data: map[string]string{
					"example": "sensitive",
				},
			}

			err := client.Create(ctx, configMap)
			Expect(err).NotTo(HaveOccurred())

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-configmap",
				},
			}

			_, err = r.Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			// Add assertions to verify the audit results
		})
	})

	Context("When creating a Secret", func() {
		It("should audit the Secret for encryption", func() {
			ctx := context.Background()

			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "default",
				},
				Data: map[string][]byte{
					"password": []byte("supersecret"),
				},
			}

			err := client.Create(ctx, secret)
			Expect(err).NotTo(HaveOccurred())

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-secret",
				},
			}

			_, err = r.Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			// Add assertions to verify the audit results
		})
	})

	Context("When creating a Role", func() {
		It("should audit the Role for excessive permissions", func() {
			ctx := context.Background()

			role := &rbacv1.Role{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-role",
					Namespace: "default",
				},
				Rules: []rbacv1.PolicyRule{
					{
						APIGroups: []string{""},
						Resources: []string{"pods"},
						Verbs:     []string{"get", "watch", "list", "delete"},
					},
				},
			}

			err := client.Create(ctx, role)
			Expect(err).NotTo(HaveOccurred())

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-role",
				},
			}

			_, err = r.Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			// Add assertions to verify the audit results
		})
	})

	Context("When creating a RoleBinding", func() {
		It("should audit the RoleBinding for role misconfigurations", func() {
			ctx := context.Background()

			roleBinding := &rbacv1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-rolebinding",
					Namespace: "default",
				},
				RoleRef: rbacv1.RoleRef{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "Role",
					Name:     "test-role",
				},
				Subjects: []rbacv1.Subject{
					{
						Kind:      "User",
						Name:      "test-user",
						Namespace: "default",
					},
				},
			}

			err := client.Create(ctx, roleBinding)
			Expect(err).NotTo(HaveOccurred())

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: "default",
					Name:      "test-rolebinding",
				},
			}

			_, err = r.Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())

			// Add assertions to verify the audit results
		})
	})
})
