#!/usr/bin/env python3
"""
Code Review Assistant
Performs systematic code reviews with security, performance, and quality checks
"""

import os
import re
import sys
import json
import yaml
import argparse
import subprocess
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from datetime import datetime
from dataclasses import dataclass, asdict
from enum import Enum

class Severity(Enum):
    CRITICAL = "critical"
    HIGH = "high"
    MEDIUM = "medium"
    LOW = "low"
    INFO = "info"

class Category(Enum):
    SECURITY = "security"
    PERFORMANCE = "performance"
    QUALITY = "quality"
    BEST_PRACTICES = "best_practices"
    STYLE = "style"

@dataclass
class Issue:
    severity: Severity
    category: Category
    file: str
    line: int
    column: int
    message: str
    recommendation: str
    code_snippet: str

class CodeReviewAssistant:
    def __init__(self, config_file: str = None):
        self.issues = []
        self.config = self._load_config(config_file)
        self.file_count = 0
        self.line_count = 0
        
    def _load_config(self, config_file: str) -> Dict:
        """Load configuration from file or use defaults"""
        default_config = {
            "enabled_checks": ["security", "performance", "quality", "best_practices"],
            "severity_threshold": "low",
            "max_file_lines": 500,
            "max_method_lines": 50,
            "max_complexity": 10,
            "ignore_paths": ["node_modules", ".git", "dist", "build", "__pycache__"],
            "file_extensions": [".py", ".js", ".ts", ".java", ".cs", ".go", ".rb", ".php", ".rs"]
        }
        
        if config_file and Path(config_file).exists():
            with open(config_file, 'r') as f:
                if config_file.endswith('.yml') or config_file.endswith('.yaml'):
                    user_config = yaml.safe_load(f)
                else:
                    user_config = json.load(f)
                default_config.update(user_config)
        
        return default_config
    
    def review_file(self, filepath: str) -> List[Issue]:
        """Review a single file"""
        path = Path(filepath)
        
        if not path.exists():
            print(f"File not found: {filepath}")
            return []
        
        # Check if file should be ignored
        for ignore_pattern in self.config["ignore_paths"]:
            if ignore_pattern in str(path):
                return []
        
        # Check file extension
        if path.suffix not in self.config["file_extensions"]:
            return []
        
        self.file_count += 1
        
        with open(path, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()
            lines = content.splitlines()
            self.line_count += len(lines)
        
        file_issues = []
        
        # Run different checks based on configuration
        if "security" in self.config["enabled_checks"]:
            file_issues.extend(self._check_security(content, filepath, lines))
        
        if "performance" in self.config["enabled_checks"]:
            file_issues.extend(self._check_performance(content, filepath, lines))
        
        if "quality" in self.config["enabled_checks"]:
            file_issues.extend(self._check_quality(content, filepath, lines))
        
        if "best_practices" in self.config["enabled_checks"]:
            file_issues.extend(self._check_best_practices(content, filepath, lines))
        
        self.issues.extend(file_issues)
        return file_issues
    
    def review_project(self, project_path: str, exclude: List[str] = None) -> List[Issue]:
        """Review entire project"""
        path = Path(project_path)
        
        if not path.exists():
            print(f"Project path not found: {project_path}")
            return []
        
        exclude = exclude or self.config["ignore_paths"]
        
        for root, dirs, files in os.walk(path):
            # Remove excluded directories from traversal
            dirs[:] = [d for d in dirs if not any(ex in d for ex in exclude)]
            
            for file in files:
                filepath = os.path.join(root, file)
                if Path(filepath).suffix in self.config["file_extensions"]:
                    self.review_file(filepath)
        
        return self.issues
    
    def _check_security(self, content: str, filepath: str, lines: List[str]) -> List[Issue]:
        """Check for security vulnerabilities"""
        issues = []
        
        # SQL Injection patterns
        sql_patterns = [
            (r'execute\([^)]*f["\'].*\{.*\}', "SQL Injection vulnerability - use parameterized queries"),
            (r'execute\([^)]*%\s*%', "SQL Injection risk - use parameterized queries"),
            (r'execute\([^)]*\+\s*[\'"]', "SQL Injection risk - avoid string concatenation"),
        ]
        
        for pattern, message in sql_patterns:
            for match in re.finditer(pattern, content):
                line_num = content[:match.start()].count('\n') + 1
                issues.append(Issue(
                    severity=Severity.CRITICAL,
                    category=Category.SECURITY,
                    file=filepath,
                    line=line_num,
                    column=match.start() - content.rfind('\n', 0, match.start()),
                    message=message,
                    recommendation="Use parameterized queries or prepared statements",
                    code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
                ))
        
        # Hardcoded secrets
        secret_patterns = [
            (r'api[_\s-]?key\s*=\s*["\'][^"\']+["\']', "Hardcoded API key detected"),
            (r'password\s*=\s*["\'][^"\']+["\']', "Hardcoded password detected"),
            (r'secret\s*=\s*["\'][^"\']+["\']', "Hardcoded secret detected"),
            (r'token\s*=\s*["\'][^"\']+["\']', "Hardcoded token detected"),
        ]
        
        for pattern, message in secret_patterns:
            for match in re.finditer(pattern, content, re.IGNORECASE):
                line_num = content[:match.start()].count('\n') + 1
                issues.append(Issue(
                    severity=Severity.CRITICAL,
                    category=Category.SECURITY,
                    file=filepath,
                    line=line_num,
                    column=match.start() - content.rfind('\n', 0, match.start()),
                    message=message,
                    recommendation="Use environment variables or secure secret management",
                    code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
                ))
        
        # XSS vulnerabilities
        xss_patterns = [
            (r'innerHTML\s*=\s*[^;]+', "XSS vulnerability - avoid innerHTML with user input"),
            (r'document\.write\([^)]*\)', "XSS risk - avoid document.write"),
            (r'eval\([^)]*\)', "Code injection risk - avoid eval"),
        ]
        
        for pattern, message in xss_patterns:
            for match in re.finditer(pattern, content):
                line_num = content[:match.start()].count('\n') + 1
                issues.append(Issue(
                    severity=Severity.HIGH,
                    category=Category.SECURITY,
                    file=filepath,
                    line=line_num,
                    column=match.start() - content.rfind('\n', 0, match.start()),
                    message=message,
                    recommendation="Sanitize user input and use safe DOM manipulation methods",
                    code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
                ))
        
        # Weak cryptography
        weak_crypto_patterns = [
            (r'hashlib\.md5\(', "Weak hashing algorithm - MD5 is broken"),
            (r'hashlib\.sha1\(', "Weak hashing algorithm - SHA1 is deprecated"),
            (r'Random\(\)', "Insecure random number generation"),
        ]
        
        for pattern, message in weak_crypto_patterns:
            for match in re.finditer(pattern, content):
                line_num = content[:match.start()].count('\n') + 1
                issues.append(Issue(
                    severity=Severity.HIGH,
                    category=Category.SECURITY,
                    file=filepath,
                    line=line_num,
                    column=match.start() - content.rfind('\n', 0, match.start()),
                    message=message,
                    recommendation="Use stronger algorithms (SHA256, bcrypt, secrets module)",
                    code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
                ))
        
        return issues
    
    def _check_performance(self, content: str, filepath: str, lines: List[str]) -> List[Issue]:
        """Check for performance issues"""
        issues = []
        
        # Nested loops (O(n²) complexity)
        nested_loop_pattern = r'for\s+.*:\s*\n\s*for\s+.*:'
        for match in re.finditer(nested_loop_pattern, content):
            line_num = content[:match.start()].count('\n') + 1
            issues.append(Issue(
                severity=Severity.MEDIUM,
                category=Category.PERFORMANCE,
                file=filepath,
                line=line_num,
                column=0,
                message="Nested loops detected - potential O(n²) complexity",
                recommendation="Consider using more efficient algorithms or data structures",
                code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
            ))
        
        # N+1 query problems (Django/SQLAlchemy)
        n_plus_one_patterns = [
            (r'for\s+\w+\s+in\s+.*\.objects\..*:\s*\n.*\.\w+\.[a-z_]+', "Potential N+1 query problem"),
            (r'for\s+\w+\s+in\s+.*\.all\(\):\s*\n.*\.\w+\.[a-z_]+', "Potential N+1 query problem"),
        ]
        
        for pattern, message in n_plus_one_patterns:
            for match in re.finditer(pattern, content):
                line_num = content[:match.start()].count('\n') + 1
                issues.append(Issue(
                    severity=Severity.HIGH,
                    category=Category.PERFORMANCE,
                    file=filepath,
                    line=line_num,
                    column=0,
                    message=message,
                    recommendation="Use select_related() or prefetch_related() for Django, joinedload() for SQLAlchemy",
                    code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
                ))
        
        # Large file operations without streaming
        file_patterns = [
            (r'open\([^)]+\)\.read\(\)', "Reading entire file into memory"),
            (r'readlines\(\)', "Loading all lines into memory"),
        ]
        
        for pattern, message in file_patterns:
            for match in re.finditer(pattern, content):
                line_num = content[:match.start()].count('\n') + 1
                issues.append(Issue(
                    severity=Severity.MEDIUM,
                    category=Category.PERFORMANCE,
                    file=filepath,
                    line=line_num,
                    column=match.start() - content.rfind('\n', 0, match.start()),
                    message=message,
                    recommendation="Use streaming or chunked reading for large files",
                    code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
                ))
        
        return issues
    
    def _check_quality(self, content: str, filepath: str, lines: List[str]) -> List[Issue]:
        """Check code quality issues"""
        issues = []
        
        # Check file length
        if len(lines) > self.config["max_file_lines"]:
            issues.append(Issue(
                severity=Severity.MEDIUM,
                category=Category.QUALITY,
                file=filepath,
                line=1,
                column=0,
                message=f"File too long ({len(lines)} lines)",
                recommendation=f"Consider splitting into smaller modules (max {self.config['max_file_lines']} lines)",
                code_snippet=""
            ))
        
        # Check method/function length
        function_pattern = r'(def |function |func |public |private |protected )[^{]+\{?'
        in_function = False
        function_start = 0
        brace_count = 0
        
        for i, line in enumerate(lines, 1):
            if re.search(function_pattern, line):
                in_function = True
                function_start = i
                brace_count = line.count('{') - line.count('}')
            elif in_function:
                brace_count += line.count('{') - line.count('}')
                if (brace_count == 0 and '{' in lines[function_start-1]) or \
                   (not '{' in lines[function_start-1] and re.match(r'^(def |function |func )', lines[function_start-1]) and not line.strip().startswith(' ')):
                    function_length = i - function_start
                    if function_length > self.config["max_method_lines"]:
                        issues.append(Issue(
                            severity=Severity.MEDIUM,
                            category=Category.QUALITY,
                            file=filepath,
                            line=function_start,
                            column=0,
                            message=f"Function too long ({function_length} lines)",
                            recommendation=f"Break down into smaller functions (max {self.config['max_method_lines']} lines)",
                            code_snippet=lines[function_start-1]
                        ))
                    in_function = False
        
        # Check for code duplication
        seen_blocks = {}
        block_size = 5
        
        for i in range(len(lines) - block_size):
            block = '\n'.join(lines[i:i+block_size])
            block_hash = hash(block)
            
            if block_hash in seen_blocks and len(block.strip()) > 50:
                issues.append(Issue(
                    severity=Severity.LOW,
                    category=Category.QUALITY,
                    file=filepath,
                    line=i+1,
                    column=0,
                    message="Duplicate code block detected",
                    recommendation="Extract duplicate code into a reusable function",
                    code_snippet=lines[i]
                ))
                break
            seen_blocks[block_hash] = i
        
        # Check for magic numbers
        magic_number_pattern = r'(?<!["\'])\b(?!0|1|2|10|100|1000)\d{2,}\b(?!["\'])'
        for i, line in enumerate(lines, 1):
            if not line.strip().startswith('#') and not line.strip().startswith('//'):
                for match in re.finditer(magic_number_pattern, line):
                    issues.append(Issue(
                        severity=Severity.LOW,
                        category=Category.QUALITY,
                        file=filepath,
                        line=i,
                        column=match.start(),
                        message=f"Magic number '{match.group()}' detected",
                        recommendation="Extract magic numbers to named constants",
                        code_snippet=line
                    ))
        
        # Check for TODO/FIXME comments
        todo_pattern = r'(TODO|FIXME|HACK|XXX)'
        for i, line in enumerate(lines, 1):
            if re.search(todo_pattern, line, re.IGNORECASE):
                issues.append(Issue(
                    severity=Severity.INFO,
                    category=Category.QUALITY,
                    file=filepath,
                    line=i,
                    column=0,
                    message="Unresolved TODO/FIXME comment",
                    recommendation="Address TODO items or create tracking issues",
                    code_snippet=line
                ))
        
        return issues
    
    def _check_best_practices(self, content: str, filepath: str, lines: List[str]) -> List[Issue]:
        """Check for best practices violations"""
        issues = []
        
        # Check for missing error handling
        try_without_except = r'try:\s*\n(?:(?!except)(?!finally).)*$'
        for match in re.finditer(try_without_except, content, re.MULTILINE):
            line_num = content[:match.start()].count('\n') + 1
            issues.append(Issue(
                severity=Severity.MEDIUM,
                category=Category.BEST_PRACTICES,
                file=filepath,
                line=line_num,
                column=0,
                message="Try block without except clause",
                recommendation="Add proper error handling",
                code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
            ))
        
        # Check for print statements in production code
        if not 'test' in filepath.lower():
            print_pattern = r'(print\(|console\.log\(|System\.out\.print)'
            for match in re.finditer(print_pattern, content):
                line_num = content[:match.start()].count('\n') + 1
                issues.append(Issue(
                    severity=Severity.LOW,
                    category=Category.BEST_PRACTICES,
                    file=filepath,
                    line=line_num,
                    column=match.start() - content.rfind('\n', 0, match.start()),
                    message="Debug print statement in production code",
                    recommendation="Use proper logging instead of print statements",
                    code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
                ))
        
        # Check for missing type hints (Python)
        if filepath.endswith('.py'):
            function_pattern = r'def\s+(\w+)\s*\(([^)]*)\)\s*:'
            for match in re.finditer(function_pattern, content):
                if '->' not in match.group():
                    line_num = content[:match.start()].count('\n') + 1
                    issues.append(Issue(
                        severity=Severity.LOW,
                        category=Category.BEST_PRACTICES,
                        file=filepath,
                        line=line_num,
                        column=0,
                        message=f"Function '{match.group(1)}' missing return type hint",
                        recommendation="Add type hints for better code documentation",
                        code_snippet=lines[line_num-1] if line_num <= len(lines) else ""
                    ))
        
        return issues
    
    def generate_report(self, format: str = "text") -> str:
        """Generate review report in specified format"""
        if format == "json":
            return self._generate_json_report()
        elif format == "html":
            return self._generate_html_report()
        elif format == "markdown":
            return self._generate_markdown_report()
        else:
            return self._generate_text_report()
    
    def _generate_text_report(self) -> str:
        """Generate plain text report"""
        report = []
        report.append("=" * 80)
        report.append("CODE REVIEW REPORT")
        report.append("=" * 80)
        report.append(f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        report.append(f"Files reviewed: {self.file_count}")
        report.append(f"Lines of code: {self.line_count}")
        report.append("")
        
        # Group issues by severity
        severity_groups = {}
        for issue in self.issues:
            if issue.severity not in severity_groups:
                severity_groups[issue.severity] = []
            severity_groups[issue.severity].append(issue)
        
        # Summary
        report.append("SUMMARY")
        report.append("-" * 40)
        total_issues = len(self.issues)
        report.append(f"Total issues: {total_issues}")
        
        for severity in Severity:
            count = len(severity_groups.get(severity, []))
            if count > 0:
                report.append(f"  {severity.value.upper()}: {count}")
        
        report.append("")
        
        # Detailed issues
        for severity in [Severity.CRITICAL, Severity.HIGH, Severity.MEDIUM, Severity.LOW, Severity.INFO]:
            issues = severity_groups.get(severity, [])
            if issues:
                report.append(f"\n{severity.value.upper()} ISSUES")
                report.append("-" * 40)
                
                for i, issue in enumerate(issues, 1):
                    report.append(f"\n{i}. {issue.message}")
                    report.append(f"   File: {issue.file}:{issue.line}")
                    report.append(f"   Category: {issue.category.value}")
                    if issue.code_snippet:
                        report.append(f"   Code: {issue.code_snippet.strip()}")
                    report.append(f"   Recommendation: {issue.recommendation}")
        
        return "\n".join(report)
    
    def _generate_json_report(self) -> str:
        """Generate JSON report"""
        report = {
            "metadata": {
                "generated": datetime.now().isoformat(),
                "files_reviewed": self.file_count,
                "lines_of_code": self.line_count,
                "total_issues": len(self.issues)
            },
            "summary": {},
            "issues": []
        }
        
        # Count by severity
        for severity in Severity:
            count = sum(1 for issue in self.issues if issue.severity == severity)
            report["summary"][severity.value] = count
        
        # Add issues
        for issue in self.issues:
            report["issues"].append({
                "severity": issue.severity.value,
                "category": issue.category.value,
                "file": issue.file,
                "line": issue.line,
                "column": issue.column,
                "message": issue.message,
                "recommendation": issue.recommendation,
                "code_snippet": issue.code_snippet
            })
        
        return json.dumps(report, indent=2)
    
    def _generate_markdown_report(self) -> str:
        """Generate Markdown report"""
        report = []
        report.append("# Code Review Report")
        report.append("")
        report.append(f"**Generated**: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        report.append(f"**Files Reviewed**: {self.file_count}")
        report.append(f"**Lines of Code**: {self.line_count}")
        report.append("")
        
        # Summary
        report.append("## Summary")
        report.append("")
        report.append(f"**Total Issues**: {len(self.issues)}")
        report.append("")
        
        severity_counts = {}
        for issue in self.issues:
            severity_counts[issue.severity] = severity_counts.get(issue.severity, 0) + 1
        
        for severity in Severity:
            count = severity_counts.get(severity, 0)
            if count > 0:
                emoji = {
                    Severity.CRITICAL: "🔴",
                    Severity.HIGH: "🟠",
                    Severity.MEDIUM: "🟡",
                    Severity.LOW: "🟢",
                    Severity.INFO: "🔵"
                }.get(severity, "")
                report.append(f"- {emoji} **{severity.value.upper()}**: {count}")
        
        report.append("")
        
        # Group issues by severity
        for severity in [Severity.CRITICAL, Severity.HIGH, Severity.MEDIUM, Severity.LOW, Severity.INFO]:
            severity_issues = [i for i in self.issues if i.severity == severity]
            if severity_issues:
                report.append(f"## {severity.value.title()} Issues")
                report.append("")
                
                for i, issue in enumerate(severity_issues, 1):
                    report.append(f"### {i}. {issue.message}")
                    report.append(f"**File**: `{issue.file}:{issue.line}`")
                    report.append(f"**Category**: {issue.category.value}")
                    report.append("")
                    
                    if issue.code_snippet:
                        report.append("```")
                        report.append(issue.code_snippet.strip())
                        report.append("```")
                        report.append("")
                    
                    report.append(f"**Recommendation**: {issue.recommendation}")
                    report.append("")
        
        return "\n".join(report)
    
    def _generate_html_report(self) -> str:
        """Generate HTML report"""
        html = []
        html.append("""<!DOCTYPE html>
<html>
<head>
    <title>Code Review Report</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            margin: 40px;
            background: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            padding: 30px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        h1 { color: #333; border-bottom: 3px solid #007acc; padding-bottom: 10px; }
        h2 { color: #555; margin-top: 30px; }
        .summary {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .critical { color: #d32f2f; font-weight: bold; }
        .high { color: #f57c00; font-weight: bold; }
        .medium { color: #fbc02d; }
        .low { color: #388e3c; }
        .info { color: #1976d2; }
        .issue {
            background: #fff;
            border-left: 4px solid;
            margin: 15px 0;
            padding: 15px;
            border-radius: 4px;
        }
        .issue.critical { border-color: #d32f2f; background: #ffebee; }
        .issue.high { border-color: #f57c00; background: #fff3e0; }
        .issue.medium { border-color: #fbc02d; background: #fffde7; }
        .issue.low { border-color: #388e3c; background: #e8f5e9; }
        .issue.info { border-color: #1976d2; background: #e3f2fd; }
        code {
            background: #f4f4f4;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }
        pre {
            background: #2d2d2d;
            color: #f8f8f2;
            padding: 15px;
            border-radius: 5px;
            overflow-x: auto;
        }
        .meta {
            color: #666;
            font-size: 0.9em;
        }
        .stats {
            display: flex;
            justify-content: space-around;
            margin: 20px 0;
        }
        .stat {
            text-align: center;
        }
        .stat-value {
            font-size: 2em;
            font-weight: bold;
            color: #007acc;
        }
        .stat-label {
            color: #666;
            font-size: 0.9em;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Code Review Report</h1>
        <p class="meta">Generated: """ + datetime.now().strftime('%Y-%m-%d %H:%M:%S') + """</p>
        """)
        
        # Statistics
        html.append(f"""
        <div class="stats">
            <div class="stat">
                <div class="stat-value">{self.file_count}</div>
                <div class="stat-label">Files Reviewed</div>
            </div>
            <div class="stat">
                <div class="stat-value">{self.line_count}</div>
                <div class="stat-label">Lines of Code</div>
            </div>
            <div class="stat">
                <div class="stat-value">{len(self.issues)}</div>
                <div class="stat-label">Total Issues</div>
            </div>
        </div>
        """)
        
        # Summary
        html.append('<div class="summary">')
        html.append('<h2>Summary</h2>')
        
        severity_counts = {}
        for issue in self.issues:
            severity_counts[issue.severity] = severity_counts.get(issue.severity, 0) + 1
        
        html.append('<ul>')
        for severity in Severity:
            count = severity_counts.get(severity, 0)
            if count > 0:
                html.append(f'<li class="{severity.value}">{severity.value.upper()}: {count}</li>')
        html.append('</ul>')
        html.append('</div>')
        
        # Issues by severity
        for severity in [Severity.CRITICAL, Severity.HIGH, Severity.MEDIUM, Severity.LOW, Severity.INFO]:
            severity_issues = [i for i in self.issues if i.severity == severity]
            if severity_issues:
                html.append(f'<h2>{severity.value.title()} Issues</h2>')
                
                for issue in severity_issues:
                    html.append(f'<div class="issue {severity.value}">')
                    html.append(f'<h3>{issue.message}</h3>')
                    html.append(f'<p><strong>File:</strong> <code>{issue.file}:{issue.line}</code></p>')
                    html.append(f'<p><strong>Category:</strong> {issue.category.value}</p>')
                    
                    if issue.code_snippet:
                        html.append('<pre><code>' + issue.code_snippet.strip() + '</code></pre>')
                    
                    html.append(f'<p><strong>Recommendation:</strong> {issue.recommendation}</p>')
                    html.append('</div>')
        
        html.append("""
    </div>
</body>
</html>""")
        
        return "\n".join(html)

def main():
    parser = argparse.ArgumentParser(description="Code Review Assistant")
    
    # Input options
    input_group = parser.add_mutually_exclusive_group(required=True)
    input_group.add_argument("--file", help="Review single file")
    input_group.add_argument("--project", help="Review entire project")
    input_group.add_argument("--git-diff", help="Review git diff (e.g., HEAD~1, main..branch)")
    
    # Check options
    parser.add_argument(
        "--checks",
        default="all",
        help="Comma-separated checks: security,performance,quality,best_practices,all"
    )
    
    # Output options
    parser.add_argument("--output", help="Output file path")
    parser.add_argument(
        "--format",
        choices=["text", "json", "html", "markdown"],
        default="text",
        help="Output format"
    )
    
    # Configuration
    parser.add_argument("--config", help="Configuration file path")
    parser.add_argument("--exclude", help="Comma-separated paths to exclude")
    
    args = parser.parse_args()
    
    # Initialize reviewer
    reviewer = CodeReviewAssistant(args.config)
    
    # Configure checks
    if args.checks != "all":
        reviewer.config["enabled_checks"] = args.checks.split(",")
    
    # Run review
    if args.file:
        reviewer.review_file(args.file)
    elif args.project:
        exclude = args.exclude.split(",") if args.exclude else None
        reviewer.review_project(args.project, exclude)
    elif args.git_diff:
        # Get changed files from git
        try:
            result = subprocess.run(
                ["git", "diff", "--name-only", args.git_diff],
                capture_output=True,
                text=True,
                check=True
            )
            files = result.stdout.strip().split("\n")
            for file in files:
                if file and Path(file).exists():
                    reviewer.review_file(file)
        except subprocess.CalledProcessError as e:
            print(f"Error getting git diff: {e}")
            sys.exit(1)
    
    # Generate report
    report = reviewer.generate_report(args.format)
    
    # Output report
    if args.output:
        with open(args.output, 'w') as f:
            f.write(report)
        print(f"Report saved to: {args.output}")
    else:
        print(report)
    
    # Exit with error code if critical issues found
    critical_count = sum(1 for i in reviewer.issues if i.severity == Severity.CRITICAL)
    if critical_count > 0:
        sys.exit(1)

if __name__ == "__main__":
    main()