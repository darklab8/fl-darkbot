import argparse


def main(args=None):
    parser = argparse.ArgumentParser(description="Process some integers.")
    parser.add_argument("option", type=str)
    parser.add_argument("abc", type=str)

    args = parser.parse_args(args)
    result = f"{repr(args)}"

    print(result)
