import uvicorn


if __name__ == "__main__":
    uvicorn.run(
        "discorder:create_app",
        host="0.0.0.0",
        factory=True,
        reload=True,
    )
