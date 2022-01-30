import os


def shell(cmd):
    print(cmd)
    status_code = os.system(cmd)

    if status_code != 0:
        exit(status_code)


shell('helm upgrade --install --create-namespace --namespace darkbot-dev darkbot . --values=darkbot_dev.yaml --values=secret_dev.yaml')
