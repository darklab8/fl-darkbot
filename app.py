import discord
import os
import json
from types import SimpleNamespace
from dotenv import load_dotenv

from discord.ext import commands

# nice settings loading
load_dotenv()
settings = SimpleNamespace()
for item, value in os.environ.items():
    setattr(settings, item, value)


class MyClient(discord.Client):
    async def on_ready(self):
        print('Logged on as {0}!'.format(self.user))

    async def on_message(self, message):
        print('Message from {0.author}: {0.content}'.format(message))

        # if message.content.startswith(command_start):
        #     if 'hello' in message.content:
        #         await message.channel.send('hi there!')
        #     if 'help' in message.content:
        #         await message.channel.send(
        #             'help command list will be implemented later!')
        # and "darkwind" in message.author.name.lower():
        #


APIDATA = SimpleNamespace()

from discord.ext import tasks, commands
import requests
import random


class MyCog(commands.Cog):
    def __init__(self, bot):
        self.index = 0
        self.bot = bot
        self.printer.start()

    def cog_unload(self):
        self.printer.cancel()
        print('unloading')

    @tasks.loop(seconds=5.0)
    async def printer(self):
        #print(self.index)
        APIDATA.players = requests.get(settings.player_request_url).json()
        APIDATA.bases = requests.get(settings.base_request_url).json()
        self.index += 1

    @printer.before_loop
    async def before_printer(self):
        print('waiting...')
        await self.bot.wait_until_ready()


bot = MyClient()

cog = MyCog(bot)

bot = commands.Bot(command_prefix='$')


@bot.command(name='fun')
async def nine_nine(ctx):
    "says random message"
    brooklyn_99_quotes = [
        'I\'m the human form of the 💯 emoji.',
        'Bingpot!',
        ('Cool. Cool cool cool cool cool cool cool, '
         'no doubt no doubt no doubt no doubt.'),
    ]

    response = random.choice(brooklyn_99_quotes)
    await ctx.send(response)


@bot.command(name='check')
async def check(ctx, number: int):
    await ctx.send(f"{number} is your lucky number!")


@check.error
async def check_error(ctx, error):
    if isinstance(error, commands.CommandError):
        await ctx.send('incorrect command!')


bot.run(settings.secret_key)