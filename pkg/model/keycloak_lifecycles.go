package model

import (
	"github.com/keycloak/keycloak-operator/pkg/apis/keycloak/v1alpha1"
	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	LifecyclePostStartImplementation = `#!/bin/bash
set -e

SCRIPT_FILE='/home/jboss/script.cli'
CHECKSUM_FILE='/home/jboss/script.checksum'

if ! [ -f ${SCRIPT_FILE} ]; then
	echo "No script file. Nothing to do."
	exit 0
fi

if [ -f ${CHECKSUM_FILE} ]; then
	/usr/bin/sha512sum -c ${CHECKSUM_FILE}
	if [ $? == 0 ]; then
	echo "Script already applied nothing to do."
	exit 0
	else
	echo "Script has changed."
	fi
fi

HTTP_STATUS=0
while [ ${HTTP_STATUS} != 200 ]
do
	HTTP_STATUS=$(curl -s -w '%{http_code}' 'http://localhost:9990/management' -o /dev/null)
	sleep 5
done

/opt/eap/bin/jboss-cli.sh --connect --file=${SCRIPT_FILE}
/usr/bin/sha512sum ${SCRIPT_FILE} > ${CHECKSUM_FILE}

#EOF
`
)

func KeycloakLifecycles(cr *v1alpha1.Keycloak) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: v12.ObjectMeta{
			Name:      KeycloakLifecyclesName,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app":           ApplicationName,
				ApplicationName: cr.Name,
			},
		},
		Data: map[string]string{
			LifecyclePostStartProperty:  LifecyclePostStartImplementation,
		},
	}
}

func KeycloakLifecyclesSelector(cr *v1alpha1.Keycloak) client.ObjectKey {
	return client.ObjectKey{
		Name:      KeycloakLifecyclesName,
		Namespace: cr.Namespace,
	}
}
