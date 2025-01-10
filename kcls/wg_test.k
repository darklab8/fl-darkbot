import infra.kcls.models.github_workflow as gw
import infra.kcls.common.github as ghc
import infra.kcls.actions.install_go as go
import infra.kcls.actions.install_taskfile as tf
import infra.kcls.models.github_action as ga

_workflow_name = "Test"
_workflow = gw.Workflow {
    _filename = "test.yml"
    name = _workflow_name
    on = {
        workflow_dispatch: {}
        push = {
            branches = ["master"]
        }
        schedule: [{cron = '0 10 1 * *'}]
    }
    jobs = {
        job = _test_job
    }
}
_test_job: gw.Job = {
    name = _workflow_name
    steps = [
        ghc.CheckoutRepo
        go.InstallGo {}
        tf.InstallTaskfile {}
        ga.Step {
            name = "test"
            run: "task test"
        }
    ]
}
