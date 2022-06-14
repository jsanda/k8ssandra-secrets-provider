package secrets

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type K8ssandraSecret struct {
	Username string
	Password string
}

type Provider interface {
	GetSecret(ctx context.Context, namespace, name string) (K8ssandraSecret, error)

	CreateSecret(ctx context.Context, namespace string, secret K8ssandraSecret) error
}

type defaultProvider struct {
	client.Client
}

func (p *defaultProvider) GetSecret(ctx context.Context, namespace, name string) (K8ssandraSecret, error) {
	k8sSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	if err := p.Get(ctx, client.ObjectKey{Namespace: namespace, Name: name}, k8sSecret); err != nil {
		return K8ssandraSecret{}, err
	}

	return K8ssandraSecret{
		Username: string(k8sSecret.Data["username"]),
		Password: string(k8sSecret.Data["password"]),
	}, nil
}

func (p *defaultProvider) CreateSecret(ctx context.Context, namespace, name string, secret K8ssandraSecret) error {
	k8sSecret := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{
			"username": []byte(secret.Username),
			"password": []byte(secret.Password),
		},
	}

	return p.Create(ctx, k8sSecret)
}