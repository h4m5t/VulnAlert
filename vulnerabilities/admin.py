# vulnerabilities/admin.py

from django.contrib import admin
from .models import Vulnerability

@admin.register(Vulnerability)
class VulnerabilityAdmin(admin.ModelAdmin):
    list_display = ('id', 'key', 'title', 'severity', 'create_time', 'update_time')
    search_fields = ('key', 'title', 'description', 'cve', 'source', 'tags')  # 使用 'source'
    list_filter = ('severity', 'create_time', 'pushed')
    ordering = ('-create_time',)
    readonly_fields = ('create_time', 'update_time')  # 设置为只读字段

    fieldsets = (
        (None, {
            'fields': (
                'key', 'title', 'description', 'severity',
                'cve', 'disclosure', 'solutions', 'references',
                'tags', 'github_search', 'source', 'pushed'
            )
        }),
        ('时间信息', {
            'fields': ('create_time', 'update_time'),
        }),
    )


# 配置说明：
# list_display: 在漏洞列表页面显示的字段，包括 id, key, title, severity, create_time, update_time。
# search_fields: 允许在管理员后台搜索的字段，包括 key, title, description, cve, source, tags。
# list_filter: 侧边栏过滤器，包括 severity, create_time, pushed。
# ordering: 默认排序方式，这里按照 create_time 逆序排列，最新的记录排在最前面。
# readonly_fields: 设置 create_time 和 update_time 为只读，防止在编辑时被修改。
# fieldsets: 组织管理员后台表单的布局，将时间信息分组显示。
# 确保 list_display 和 list_filter 中引用的字段名称与 Vulnerability 模型中的字段名称完全一致。