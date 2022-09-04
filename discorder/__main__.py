import uvicorn


if __name__ == "__main__":
    uvicorn.run("discorder:create_app", factory=True, reload=True)
