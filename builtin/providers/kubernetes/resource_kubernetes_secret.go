package kubernetes

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/kubernetes/pkg/api/errors"
	api "k8s.io/kubernetes/pkg/api/v1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func resourceKubernetesSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesSecretCreate,
		Read:   resourceKubernetesSecretRead,
		Exists: resourceKubernetesSecretExists,
		Update: resourceKubernetesSecretUpdate,
		Delete: resourceKubernetesSecretDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("secret", true),
			"data": {
				Type:        schema.TypeMap,
				Description: "A map of the secret data.",
				Optional:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of secret",
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceKubernetesSecretCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	secret := api.Secret{
		ObjectMeta: metadata,
		StringData: expandStringMap(d.Get("data").(map[string]interface{})),
	}

	if v, ok := d.GetOk("type"); ok {
		secret.Type = api.SecretType(v.(string))
	}

	log.Printf("[INFO] Creating new secret: %#v", secret)
	out, err := conn.CoreV1().Secrets(metadata.Namespace).Create(&secret)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Submitting new secret: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesSecretRead(d, meta)
}

func resourceKubernetesSecretRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Reading secret %s", name)
	secret, err := conn.CoreV1().Secrets(namespace).Get(name)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Received secret: %#v", secret)
	err = d.Set("metadata", flattenMetadata(secret.ObjectMeta))
	if err != nil {
		return err
	}

	data := map[string]string{}
	for k, v := range secret.Data {
		data[k] = string(v)
	}

	d.Set("data", data)
	d.Set("type", secret.Type)

	return nil
}

func resourceKubernetesSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	namespace, name := idParts(d.Id())
	metadata.Name = name

	secret := api.Secret{
		ObjectMeta: metadata,
		StringData: expandStringMap(d.Get("data").(map[string]interface{})),
	}

	if v, ok := d.GetOk("type"); ok {
		secret.Type = api.SecretType(v.(string))
	}

	log.Printf("[INFO] Updating secret: %#v", secret)
	out, err := conn.CoreV1().Secrets(namespace).Update(&secret)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Submitting updated secret: %#v", out)
	d.SetId(buildId(out.ObjectMeta))

	return resourceKubernetesSecretRead(d, meta)
}

func resourceKubernetesSecretDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Deleting secret: %#v", name)
	err := conn.CoreV1().Secrets(namespace).Delete(name, &api.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Secret %s deleted", name)

	d.SetId("")

	return nil
}

func resourceKubernetesSecretExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*kubernetes.Clientset)

	namespace, name := idParts(d.Id())

	log.Printf("[INFO] Checking secret %s", name)
	_, err := conn.CoreV1().Secrets(namespace).Get(name)
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			return false, nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
	}

	return true, err
}
