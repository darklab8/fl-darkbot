import discord
from .connector import run_command_async


class MyClient(discord.Client):
    async def on_ready(self):
        print("Logged on as", self.user)

    async def on_message(self, message):
        # don't respond to ourselves
        if message.author == self.user:
            return

        content: str = message.content
        print(f"content={content}")

        prefix = ".bot"
        if content.startswith(prefix):
            content = content[len(prefix) :]

            args = [arg for arg in content.split(" ") if arg != ""]
            print(f"args={args}")

            result = await run_command_async(*(["python3", "-m" "consoler"] + args))

            # print(f"answer={result}")

            # await message.channel.send("pong")
            await message.channel.send(result)
