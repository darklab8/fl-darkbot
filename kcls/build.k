import .wg_test
import .wg_pages
import infra.kcls.common.workflow
import yaml

_image_name = "darkwind8/darkbot"

_workflows = [
    workflow.DockerDeploySwarm(workflow.DockerDeploySwarmOpts {
        image_name = _image_name
        service_name = "darkbot"
        environment = "staging"
    })
    workflow.DockerDeploySwarm(workflow.DockerDeploySwarmOpts {
        image_name = _image_name
        service_name = "darkbot"
        environment = "production"
    })
    wg_test._workflow
    wg_pages._workflow
]

_run = [yaml.dump_to_file(a, ".github/workflows/${a._filename}", ignore_private=True, ignore_none=True) for a in _workflows]
