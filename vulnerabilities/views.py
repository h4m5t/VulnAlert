from django.shortcuts import render, get_object_or_404
from django.db.models import Count, Q
from django.db.models.functions import TruncMonth, TruncDay
from django.contrib.auth.decorators import login_required
from django.utils import timezone
from datetime import datetime, timedelta
from .models import Vulnerability
import json
from django.core.paginator import Paginator, EmptyPage, PageNotAnInteger
import re

@login_required
def home(request):
    """
    系统概览页面视图
    显示关键统计数据、趋势图表和最新漏洞列表
    """
    # 获取当前时间和时间范围
    now = timezone.now()
    start_of_month = now.replace(day=1, hour=0, minute=0, second=0, microsecond=0)
    start_of_day = now.replace(hour=0, minute=0, second=0, microsecond=0)  # 添加这行
    six_months_ago = now - timedelta(days=180)


    # 1. 基础统计数据
    total_vulnerabilities = Vulnerability.objects.count()
    
    # 2. 严重程度统计
    severity_counts = dict(
        Vulnerability.objects.values('severity')
        .annotate(count=Count('id'))
        .values_list('severity', 'count')
    )
    
    # 计算严重漏洞百分比
    critical_count = severity_counts.get('严重', 0)
    severity_percentage = round((critical_count / total_vulnerabilities * 100) if total_vulnerabilities > 0 else 0, 1)

    # 3. 本月新增统计
    monthly_new = Vulnerability.objects.filter(
        create_time__gte=start_of_month
    ).count()


        # 添加本日新增统计
    daily_new = Vulnerability.objects.filter(
        create_time__gte=start_of_day
    ).count()

    # 4. 已推送统计
    pushed_count = Vulnerability.objects.filter(pushed=True).count()

    # 5. 趋势数据处理
    trend_labels = []
    trend_counts = []
    
    # 生成最近6个月的月份列表
    for i in range(5, -1, -1):
        # 计算月份
        current_date = now - timedelta(days=30 * i)
        month_start = current_date.replace(day=1, hour=0, minute=0, second=0, microsecond=0)
        month_end = (month_start + timedelta(days=32)).replace(day=1) - timedelta(seconds=1)
        
        # 获取该月的漏洞数量
        month_count = Vulnerability.objects.filter(
            create_time__gte=month_start,
            create_time__lte=month_end
        ).count()
        
        # 添加到列表
        trend_labels.append(month_start.strftime('%Y-%m'))
        trend_counts.append(month_count)

    # 6. 按漏洞来源统计
    source_stats = (
        Vulnerability.objects.values('source')
        .annotate(count=Count('id'))
        .order_by('-count')
    )

    # 7. 最新漏洞列表
    latest_vulnerabilities = (
        Vulnerability.objects.all()
        .order_by('-create_time')
        [:10]  # 限制显示最新的10条记录
    )

    # 8. 计算同比增长率
    last_month = now - timedelta(days=30)
    last_month_count = Vulnerability.objects.filter(
        create_time__gte=last_month
    ).count()
    
    previous_month = last_month - timedelta(days=30)
    previous_month_count = Vulnerability.objects.filter(
        create_time__gte=previous_month,
        create_time__lt=last_month
    ).count()
    
    growth_rate = (
        ((last_month_count - previous_month_count) / previous_month_count * 100)
        if previous_month_count > 0 else 0
    )

    # 构建上下文数据
    context = {
        # 基础统计
        'total_vulnerabilities': total_vulnerabilities,
        'severity_counts': severity_counts,
        'monthly_new': monthly_new,
        'pushed_count': pushed_count,
        'severity_percentage': severity_percentage,
        
        # 趋势数据
        'trend_labels': json.dumps(trend_labels),  # 转换为JSON字符串
        'trend_data': trend_counts,
        
        # 来源统计
        'source_stats': source_stats,
        
        # 最新漏洞
        'latest_vulnerabilities': latest_vulnerabilities,
        
        # 增长统计
        'growth_rate': round(growth_rate, 1),
        'daily_new': daily_new,  # 添加这行
        # 页面标题
        'page_title': '系统概览',
    }

    return render(request, 'vulnerabilities/home.html', context)

@login_required
def vulnerability_list(request):
    # 获取所有漏洞，按创建时间降序排列，排除严重性为空的漏洞
    vulnerabilities = Vulnerability.objects.all().order_by('-create_time').exclude(severity__isnull=True).exclude(severity='')
    
    # 获取搜索参数
    query = request.GET.get('q', '').strip()
    
    # 获取严重性过滤参数
    severity_filter = request.GET.get('severity', '').strip()

    pushed_filter = request.GET.get('pushed', '')  # 新增获取推送状态参数
    
    # 应用搜索条件
    if query:
        vulnerabilities = vulnerabilities.filter(
            Q(key__icontains=query) |
            Q(cve__icontains=query) |
            Q(title__icontains=query) |
            Q(description__icontains=query)  # 新增描述字段的搜索
        )
    
    # 应用严重性过滤条件
    if severity_filter and severity_filter in dict(Vulnerability.SEVERITY_CHOICES):
        vulnerabilities = vulnerabilities.filter(severity=severity_filter)
    
    # 处理推送状态过滤
    if pushed_filter:
        if pushed_filter.lower() in ['true', '1', 'yes', '已推送']:
            vulnerabilities = vulnerabilities.filter(pushed=True)
        elif pushed_filter.lower() in ['false', '0', 'no', '未推送']:
            vulnerabilities = vulnerabilities.filter(pushed=False)
    
    # 分页设置
    paginator = Paginator(vulnerabilities, 15)  # 每页显示 15 条
    page = request.GET.get('page')
    
    try:
        vulnerabilities_page = paginator.page(page)
    except PageNotAnInteger:
        vulnerabilities_page = paginator.page(1)
    except EmptyPage:
        vulnerabilities_page = paginator.page(paginator.num_pages)
    
    context = {
        'vulnerabilities': vulnerabilities_page,
        'query': query,
        'severity_filter': severity_filter,
        'pushed_filter': pushed_filter,  # 确保推送状态过滤参数传递到模板
        'severity_choices': Vulnerability.SEVERITY_CHOICES,
    }
    return render(request, 'vulnerabilities/vulnerability_list.html', context)

@login_required
def vulnerability_detail(request, pk):
    vulnerability = get_object_or_404(Vulnerability, pk=pk)
    
    # 处理references
    parsed_references = []
    if vulnerability.references:
        if isinstance(vulnerability.references, list):
            # 直接使用整个列表
            parsed_references = vulnerability.references
        elif isinstance(vulnerability.references, str):
            # 如果是字符串，先尝试解析JSON
            try:
                import json
                parsed_references = json.loads(vulnerability.references)
            except json.JSONDecodeError:
                # 如果JSON解析失败，则按分隔符分割
                parsed_references = [ref.strip() for ref in re.split(r'[;,]', vulnerability.references) if ref.strip()]
    
    context = {
        'vulnerability': vulnerability,
        'parsed_references': parsed_references
    }
    return render(request, 'vulnerabilities/vulnerability_detail.html', context)