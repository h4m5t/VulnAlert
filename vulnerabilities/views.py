# vulnerabilities/views.py

import re
from django.shortcuts import render, get_object_or_404
from .models import Vulnerability
from django.db import models

def home(request):
    total_vulnerabilities = Vulnerability.objects.count()
    severity_counts = Vulnerability.objects.values('severity').annotate(count=models.Count('severity'))
    # 转换为字典
    severity_counts_dict = {item['severity']: item['count'] for item in severity_counts}
    context = {
        'total_vulnerabilities': total_vulnerabilities,
        'severity_counts': severity_counts_dict,
    }
    return render(request, 'vulnerabilities/home.html', context)

def vulnerability_list(request):
    vulnerabilities = Vulnerability.objects.all()
    context = {
        'vulnerabilities': vulnerabilities
    }
    return render(request, 'vulnerabilities/vulnerability_list.html', context)

def vulnerability_detail(request, pk):
    vulnerability = get_object_or_404(Vulnerability, pk=pk)
    references = vulnerability.references

    parsed_references = []
    if references and isinstance(references, list):
        refs_str = references[0]
        parsed_references = [ref.strip() for ref in re.split(r'[;,]', refs_str) if ref.strip()]
    
    context = {
        'vulnerability': vulnerability,
        'parsed_references': parsed_references
    }
    return render(request, 'vulnerabilities/vulnerability_detail.html', context)
