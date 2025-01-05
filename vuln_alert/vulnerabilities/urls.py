# vulnerabilities/urls.py

from django.urls import path
from . import views

app_name = 'vulnerabilities'

urlpatterns = [
    path('', views.home, name='home'),
    path('vulnerabilities/', views.vulnerability_list, name='vulnerability_list'),
    path('vulnerabilities/<int:pk>/', views.vulnerability_detail, name='vulnerability_detail'),
    # 将来可以添加创建、编辑、删除漏洞的路由
]