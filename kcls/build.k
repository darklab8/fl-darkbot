import .wf_docker_staging
import .wf_docker_production
import .wg_test
import .wg_pages

import yaml

_workflows = [
    wf_docker_staging._workflow
    wf_docker_production._workflow
    wg_test._workflow
    wg_pages._workflow
]

_run = [yaml.dump_to_file(a, ".github/workflows/${a._filename}", ignore_private=True, ignore_none=True) for a in _workflows]
