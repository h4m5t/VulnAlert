# vulnerabilities/models.py

from django.db import models

class Vulnerability(models.Model):
    SEVERITY_CHOICES = [
        ('低', '低'),
        ('中', '中'),
        ('高', '高'),
        ('严重', '严重'),
    ]

    id = models.IntegerField(primary_key=True)
    key = models.CharField(max_length=100, unique=True, help_text="唯一标识符")
    title = models.CharField(max_length=255, help_text="漏洞标题")
    description = models.TextField(help_text="漏洞描述")
    severity = models.CharField(max_length=10, choices=SEVERITY_CHOICES, help_text="漏洞严重性")
    cve = models.CharField(max_length=50, blank=True, null=True, help_text="CVE编号")
    disclosure = models.DateField(blank=True, null=True, help_text="披露日期")
    solutions = models.TextField(blank=True, null=True, help_text="解决方案")
    references = models.JSONField(blank=True, null=True, help_text="参考链接列表")
    tags = models.JSONField(blank=True, null=True, help_text="标签，逗号分隔")
    github_search = models.JSONField(blank=True, null=True, help_text="GitHub搜索链接")
    source = models.CharField(
        max_length=255,
        db_column='from',  # 映射到数据库的 'from' 列，由于 from 是 Python 的保留关键字
        help_text="漏洞来源"
    )
    pushed = models.BooleanField(default=True, help_text="是否已推送")
    create_time = models.DateTimeField(auto_now_add=True, help_text="创建时间")
    update_time = models.DateTimeField(auto_now=True, help_text="更新时间")

    class Meta:
        db_table = 'vuln_informations'  # 指定数据库表名
        managed = False  # Django 不管理此表
        verbose_name = "漏洞信息"
        verbose_name_plural = "漏洞信息"

    def __str__(self):
        return self.title