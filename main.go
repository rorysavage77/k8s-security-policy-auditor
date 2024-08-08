package main

import (
	"context"
	"os"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func main() {
	// Set up the logger
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
	log := ctrl.Log.WithName("k8s-security-policy-auditor")

	// Create a new manager to provide shared dependencies and start components
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: runtime.NewScheme(),
	})
	if err != nil {
		log.Error(err, "unable to set up overall controller manager")
		os.Exit(1)
	}

	// Add the core and RBAC APIs to the manager's scheme
	if err := corev1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "unable to add corev1 scheme")
		os.Exit(1)
	}
	if err := rbacv1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "unable to add rbacv1 scheme")
		os.Exit(1)
	}

	// Setup the controller with manager
	if err := add(mgr); err != nil {
		log.Error(err, "unable to create controller", "controller", "SecurityPolicyAuditor")
		os.Exit(1)
	}

	// Start the manager
	log.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to run manager")
		os.Exit(1)
	}
}

// Add sets up the controller with the Manager
func add(mgr manager.Manager) error {
	r := &reconciler{
		client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
	}

	// Create a new controller
	c, err := controller.New("k8s-security-policy-auditor", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch ConfigMaps
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.ConfigMap{}), &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// Watch Secrets
	if err := c.Watch(source.Kind(mgr.GetCache(), &corev1.Secret{}), &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// Watch RBAC Roles
	if err := c.Watch(source.Kind(mgr.GetCache(), &rbacv1.Role{}), &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// Watch RoleBindings
	if err := c.Watch(source.Kind(mgr.GetCache(), &rbacv1.RoleBinding{}), &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	return nil
}

// reconciler reconciles a ConfigMap, Secret, Role, or RoleBinding object
type reconciler struct {
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads the state of the cluster for a ConfigMap, Secret, Role, or RoleBinding object and makes changes based on the SecurityPolicyAuditor
func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.Log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)

	// Log the event
	log.Info("Reconciling resource", "Namespace", req.Namespace, "Name", req.Name)

	// Perform the security check
	if err := r.performSecurityAudit(ctx, req); err != nil {
		log.Error(err, "Security audit failed")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// performSecurityAudit checks the resource for security best practices
func (r *reconciler) performSecurityAudit(ctx context.Context, req ctrl.Request) error {
	// Check ConfigMaps for sensitive data
	configMap := &corev1.ConfigMap{}
	if err := r.client.Get(ctx, req.NamespacedName, configMap); err == nil {
		for key, value := range configMap.Data {
			if isSensitive(key, value) {
				ctrl.Log.Info("Sensitive data found in ConfigMap", "Namespace", configMap.Namespace, "Name", configMap.Name)
			}
		}
	}

	// Check Secrets for encryption
	secret := &corev1.Secret{}
	if err := r.client.Get(ctx, req.NamespacedName, secret); err == nil {
		// Placeholder for Secret audit logic
		ctrl.Log.Info("Auditing Secret", "Namespace", secret.Namespace, "Name", secret.Name)
	}

	// Check RBAC roles for excessive permissions
	role := &rbacv1.Role{}
	if err := r.client.Get(ctx, req.NamespacedName, role); err == nil {
		for _, rule := range role.Rules {
			if hasExcessivePermissions(rule) {
				ctrl.Log.Info("Role has excessive permissions", "Namespace", role.Namespace, "Name", role.Name)
			}
		}
	}

	// Check RoleBindings for role misconfigurations
	roleBinding := &rbacv1.RoleBinding{}
	if err := r.client.Get(ctx, req.NamespacedName, roleBinding); err == nil {
		// Placeholder for RoleBinding audit logic
		ctrl.Log.Info("Auditing RoleBinding", "Namespace", roleBinding.Namespace, "Name", roleBinding.Name)
	}

	return nil
}

// isSensitive checks if the data is considered sensitive
func isSensitive(key, value string) bool {
	// Placeholder for sensitive data detection logic
	return false
}

// hasExcessivePermissions checks if the role has excessive permissions
func hasExcessivePermissions(rule rbacv1.PolicyRule) bool {
	// Placeholder for excessive permissions detection logic
	return false
}
