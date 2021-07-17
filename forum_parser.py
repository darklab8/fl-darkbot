from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import bs4
from dataclasses import dataclass


@dataclass
class forum_record:
    title: str
    thread_author: str
    url: str
    category: str
    date: str
    views: str
    last_author: str
    replies: str
    index: int


def get_forum_threads(forum_acc: str, forum_pass: str):

    chrome_options = Options()
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--disable-dev-shm-usage")
    chrome_options.add_argument("--disable-gpu")
    chrome_options.add_argument("--headless")
    driver = webdriver.Chrome(executable_path="./chromedriver",
                              options=chrome_options)

    driver.get("https://discoverygc.com/forums/search.php?action=getdaily")
    login_href_xpath = '//a[@href="member.php?action=login"]'
    WebDriverWait(driver, 10).until(
        EC.presence_of_element_located((By.XPATH, login_href_xpath)))
    elem = driver.find_element_by_xpath(login_href_xpath)
    elem.click()

    input_username_xpath = '//input[@name="username"]'
    WebDriverWait(driver, 10).until(
        EC.presence_of_element_located((By.XPATH, input_username_xpath)))
    input1 = driver.find_element_by_xpath(input_username_xpath)
    input2 = driver.find_element_by_xpath('//input[@name="password"]')

    input1.send_keys(forum_acc)
    input2.send_keys(forum_pass, Keys.ENTER)

    result_xpath = '//tr[@class="inline_row"]'
    WebDriverWait(driver, 10).until(
        EC.presence_of_element_located((By.XPATH, result_xpath)))

    with open("resp.html", "w") as file_:
        file_.write(driver.page_source)

    bs = bs4.BeautifulSoup(driver.page_source, "html.parser")

    results = bs.find_all("tr", {"class": "inline_row"})

    items = []
    for index, item in enumerate(results):
        items.append(
            forum_record(
                title=next(
                    title for title in item.find_all("a")
                    if title.attrs.get("class") is not None
                    and "subject" in title.attrs.get("class")[0]).get_text(),
                thread_author=item.find("div", {
                    "class": "author"
                }).get_text(),
                url="https://discoverygc.com/forums/" +
                item.find_all("td")[6].a.attrs['href'],
                category=item.find_all("td")[3].a.get_text(),
                date=item.find_all("td")[6].span.span.attrs['title'],
                views=item.find_all("td")[5].get_text(),
                last_author=item.find_all("td")[6].find_all("a")[1].get_text(),
                replies=item.find_all("td")[4].get_text(),
                index=index,
            ))

    driver.close()
    return items
