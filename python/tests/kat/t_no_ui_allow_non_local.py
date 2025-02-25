from kat.harness import Query, EDGE_STACK

from abstract_tests import AmbassadorTest, ServiceType, HTTP


class NoUITestAllowNoLocal(AmbassadorTest):
    # Don't use single_namespace -- we want CRDs, so we want
    # the cluster-scope RBAC instead of the namespace-scope
    # RBAC. Our ambassador_id filters out the stuff we want.
    namespace = "no-ui-namespace-anl"
    extra_ports = [8877]

    def manifests(self) -> str:
        return (
            self.format(
                """
---
apiVersion: v1
kind: Namespace
metadata:
  name: no-ui-namespace
---
apiVersion: getambassador.io/v3alpha1
kind: Module
metadata:
  name: ambassador
  namespace: no-ui-namespace
  labels:
    kat-ambassador-id: {self.ambassador_id}
spec:
  ambassador_id: [ {self.ambassador_id} ]
  config:
    diagnostics:
      enabled: false
      allow_non_local: true
"""
            )
            + super().manifests()
        )

    def queries(self):
        yield (Query(self.url("ambassador/v0/diag/"), expected=404))
        yield (Query(self.url("edge_stack/admin/"), expected=404))
        yield Query(self.url("ambassador/v0/diag/", scheme="http", port=8877), expected=200)
