import discord
from .connector import run_command_async
from listener.permissions import PermissionChecker
import inspect


class MyClient(discord.Client):
    async def on_ready(self):
        print("Logged on as", self.user)

    async def on_message(self, message):
        # don't respond to ourselves
        if message.author == self.user:
            return

        content: str = message.content
        print(f"content={content}")

        permissions = PermissionChecker(message)
        print(f"permissions.all={permissions.is_having_any_permission}")
        print(f"permissions.predicate_controller={permissions.predicate_controller}")
        print(f"permissions.predicate_guild_owner={permissions.predicate_guild_owner}")
        print(f"permissions.predicate_manage={permissions.predicate_manage}")
        print(f"permissions.predicate_owner={permissions.predicate_owner}")

        if not permissions.is_having_any_permission:
            await message.channel.send(
                "not authorized user. required to have role 'bot_controller', or right to `manage channels`, or being guild owner"
            )
            return

        prefix = ".bot"
        if not content.startswith(prefix):
            return

        content = content[len(prefix) :]

        args = [arg for arg in content.split(" ") if arg != ""]
        print(f"args={args}")
        final_args = (
            ["python3", "-m", "consoler"]
            + args
            + [f"--channel_id={message.channel.id}"]
        )
        print(f"final_args={final_args}")
        result = await run_command_async(*final_args)

        markdown_template = inspect.cleandoc(
            """```
        {text}
        ```"""
        )
        await message.channel.send(markdown_template.format(text=result))
