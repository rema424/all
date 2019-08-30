# -*- coding: utf-8 -*-
import scrapy
from scrapy.linkextractors import LinkExtractor
from scrapy.spiders import CrawlSpider, Rule
from cartune.items import PopularCarArticleItem


class PopularCarCrawlSpider(CrawlSpider):
    name = 'popular_car_crawl'
    allowed_domains = ['cartune.me']
    start_urls = ['http://cartune.me/']

    rules = (
        Rule(
            LinkExtractor(
                restrict_css='div.content-left ul.menuList--thumbnail a.title'),
            callback='parse_item',
            follow=True,
        ),
    )

    def parse_item(self, response):
        car = response.css(
            'h1.categoryHeader-title::text').extract_first()
        for article in response.css('div.articleList--s div.item'):
            item = PopularCarArticleItem()
            item['car'] = car
            item['article_title'] = article.css(
                'div.item-title a::text').extract_first()
            item['article_url'] = article.css(
                'div.item-title a::attr(href)').extract_first()
            yield item

        # item = {}
        #item['domain_id'] = response.xpath('//input[@id="sid"]/@value').get()
        #item['name'] = response.xpath('//div[@id="name"]').get()
        #item['description'] = response.xpath('//div[@id="description"]').get()
        # return item
