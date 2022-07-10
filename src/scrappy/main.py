from fastapi import FastAPI

from .players import views

app = FastAPI()


@app.get("/")
async def get_ping():
    return {"message": "pong!"}


app.include_router(views.router)
