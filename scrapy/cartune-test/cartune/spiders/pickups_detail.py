# -*- coding: utf-8 -*-
import scrapy


class PickupsDetailSpider(scrapy.Spider):
    name = 'pickups_detail'
    allowed_domains = ['cartune.me']
    start_urls = ['http://cartune.me/']

    def parse(self, response):
        pass
