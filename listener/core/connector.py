import subprocess
import asyncio

def run_command_sync(args=[]):
    result = subprocess.run(
        ["python3", "-m" "consoler"] + args,
        capture_output=True,
        text=True,
    )
    if result.returncode != 0:
        return result.stderr
    return result.stdout

async def run_command_async(*args):
    """Run command in subprocess.

    Example from:
        http://asyncio.readthedocs.io/en/latest/subprocess.html
    """
    process = await asyncio.create_subprocess_exec(
        *args, stdout=asyncio.subprocess.PIPE, stderr=asyncio.subprocess.PIPE
    )

    stdout, stderr = await process.communicate()
    return stdout.decode() + stderr.decode()

if __name__ == "__main__":
    # print(run_command_sync(args=["ping"]))
    print(asyncio.run(run_command_async(*(["python3", "-m" "consoler"] + ["ping"]))))
