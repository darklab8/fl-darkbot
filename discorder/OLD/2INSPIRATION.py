# shared code to you how to have Discord with FastAPI. Implement in your code.

import asyncio

import discord
import uvicorn
from discord.ext import commands
from fastapi import FastAPI
from starlette.requests import Request


async def create_bot():
    intents = discord.Intents.default()
    intents.message_content = True
    intents.members = True

    bot = commands.Bot(
        command_prefix="!",
        intents=intents,
        case_insensitive=True,
    )

    @bot.command("test")
    async def test(ctx: commands.Context):
        await ctx.send("Test!")

    return bot


async def run_bot(bot: commands.Bot):
    async with bot:
        await bot.start(token="")


def create_app() -> FastAPI:
    app = FastAPI()

    @app.on_event("startup")
    async def on_startup():
        app.discord_bot = await create_bot()
        asyncio.create_task(run_bot(app.discord_bot))

    @app.get("/")
    def root(request: Request):
        return repr(request.app.discord_bot)

    return app


if __name__ == "__main__":
    uvicorn.run("test_dpy:create_app", factory=True, reload=True)
