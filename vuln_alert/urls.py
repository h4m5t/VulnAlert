# vuln_alert/urls.py

from django.contrib import admin
from django.urls import path, include
from django.contrib.auth import views as auth_views

urlpatterns = [
    path('admin/', admin.site.urls),
    path('vulnerabilities/', include('vulnerabilities.urls')),  # 包含应用的 URL 模式
    path('login/', auth_views.LoginView.as_view(template_name='vulnerabilities/login.html'), name='login'),  # 登录 URL
    path('logout/', auth_views.LogoutView.as_view(next_page='login'), name='logout'),  # 注销 URL
    path('', include('vulnerabilities.urls')),  # 让应用处理根 URL
]