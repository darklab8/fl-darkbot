import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.models.github_action as ga

_workflow_name = "Pages"
_workflow = gw.Workflow {
    _filename = "pages.yml"
    name = _workflow_name
    on = {
        workflow_dispatch: {}
        push = {
            branches = ["master"]
        }
    }
    jobs = {
        job = _pages_job
    }
}
_pages_job: gw.Job = {
    name = _workflow_name
    permissions: ghc.PagesPermissions {}
    environment: dict(ghc.PagesEnvironment {})
    steps = [
        ghc.CheckoutRepo
        ga.Step {
            uses = "actions/setup-python@v4"
            with = {"python-version": "3.10"}
        }
        ga.Step {
            name = "ensure pip"
            run = "python3 -m ensurepip"
        }
        ga.Step {
            name = "Install requirements.txt"
            run = "pip install -r requirements.txt -c constraints.txt"
        }
        ga.Step {
            name = "Build docs"
            run = "mkdocs build"
        }
    ] + ghc.PagesToDeploySteps(ghc.PagesDeployOpts {
        path = "./site"
    })
}
