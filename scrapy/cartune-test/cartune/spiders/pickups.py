# -*- coding: utf-8 -*-
import scrapy
from cartune.items import CartuneItem


class PickupsSpider(scrapy.Spider):
    name = 'pickups'
    allowed_domains = ['cartune.me']
    start_urls = ['http://cartune.me/']

    def parse(self, response):
        # pass
        for pickup in response.css("div.pickupList div.pickupList-item"):
            item = CartuneItem()
            item["car"] = pickup.css("p a::text").extract_first()
            item["url"] = pickup.css("a::attr(href)").extract_first()
            yield item
