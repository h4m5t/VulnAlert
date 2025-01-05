# vulnerabilities/urls.py

# from django.urls import path
# from .views import HomeView, VulnerabilityListView, VulnerabilityDetailView

# app_name = 'vulnerabilities'

# urlpatterns = [
#     path('', HomeView.as_view(), name='home'),  # 首页视图
#     path('list/', VulnerabilityListView.as_view(), name='vulnerability_list'),  # 漏洞列表视图
#     path('detail/<int:pk>/', VulnerabilityDetailView.as_view(), name='vulnerability_detail'),  # 漏洞详情视图
#     # 其他 URL 模式...
# ]

# vulnerabilities/urls.py

from django.urls import path
from . import views

app_name = 'vulnerabilities'

urlpatterns = [
    path('', views.home, name='home'),  # 首页视图
    path('list/', views.vulnerability_list, name='vulnerability_list'),  # 漏洞列表视图
    path('detail/<int:pk>/', views.vulnerability_detail, name='vulnerability_detail'),  # 漏洞详情视图
    # 其他 URL 模式...
]