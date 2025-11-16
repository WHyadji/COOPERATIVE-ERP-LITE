#!/usr/bin/env python3
"""
Code Style Enforcer
Apply consistent formatting, linting, and naming conventions across languages
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
from dataclasses import dataclass, asdict

@dataclass
class StyleIssue:
    file: str
    line: int
    column: int
    rule: str
    severity: str
    message: str
    fixed: bool = False
    suggestion: str = ""

class CodeStyleEnforcer:
    def __init__(self, config_file: str = None):
        self.config = self._load_config(config_file)
        self.issues = []
        self.files_checked = 0
        self.files_fixed = 0
        
    def _load_config(self, config_file: str) -> Dict:
        """Load configuration from file or use defaults"""
        default_config = {
            "python": {
                "indent_size": 4,
                "max_line_length": 88,
                "style_guide": "pep8"
            },
            "javascript": {
                "indent_size": 2,
                "max_line_length": 80,
                "style_guide": "airbnb"
            },
            "ignore_paths": ["node_modules", ".git", "dist", "build", "__pycache__"]
        }
        
        if config_file and Path(config_file).exists():
            with open(config_file, 'r') as f:
                if config_file.endswith('.yml') or config_file.endswith('.yaml'):
                    user_config = yaml.safe_load(f)
                else:
                    user_config = json.load(f)
                default_config.update(user_config)
        
        return default_config
    
    def check_file(self, filepath: str, fix: bool = False) -> List[StyleIssue]:
        """Check and optionally fix style issues in a file"""
        if not Path(filepath).exists():
            print(f"File not found: {filepath}")
            return []
        
        self.files_checked += 1
        
        with open(filepath, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()
            original_content = content
        
        file_issues = []
        
        # Detect language and apply appropriate checks
        if filepath.endswith('.py'):
            issues, fixed_content = self._check_python(content, filepath, fix)
        elif filepath.endswith(('.js', '.jsx', '.ts', '.tsx')):
            issues, fixed_content = self._check_javascript(content, filepath, fix)
        else:
            issues, fixed_content = self._check_universal(content, filepath, fix)
        
        file_issues.extend(issues)
        
        # Write fixed content if changes were made
        if fix and fixed_content != original_content:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(fixed_content)
            self.files_fixed += 1
            for issue in file_issues:
                issue.fixed = True
        
        self.issues.extend(file_issues)
        return file_issues
    
    def _check_python(self, content: str, filepath: str, fix: bool) -> Tuple[List[StyleIssue], str]:
        """Check Python style issues"""
        issues = []
        lines = content.split('\n')
        fixed_lines = []
        
        for i, line in enumerate(lines, 1):
            fixed_line = line
            
            # Check indentation
            if line and line[0] in ' \t':
                indent = len(line) - len(line.lstrip())
                if '\t' in line[:indent]:
                    issues.append(StyleIssue(
                        file=filepath,
                        line=i,
                        column=1,
                        rule="indentation",
                        severity="error",
                        message="Tabs used for indentation",
                        suggestion="Use 4 spaces"
                    ))
                    if fix:
                        fixed_line = line.replace('\t', '    ')
                elif indent % 4 != 0:
                    issues.append(StyleIssue(
                        file=filepath,
                        line=i,
                        column=1,
                        rule="indentation",
                        severity="warning",
                        message=f"Indentation not multiple of 4 (found {indent})",
                        suggestion=f"Use {(indent // 4) * 4} spaces"
                    ))
            
            # Check line length
            if len(line) > self.config["python"]["max_line_length"]:
                issues.append(StyleIssue(
                    file=filepath,
                    line=i,
                    column=self.config["python"]["max_line_length"] + 1,
                    rule="line-length",
                    severity="warning",
                    message=f"Line too long ({len(line)} > {self.config['python']['max_line_length']})",
                    suggestion="Break into multiple lines"
                ))
            
            # Check trailing whitespace
            if line != line.rstrip():
                issues.append(StyleIssue(
                    file=filepath,
                    line=i,
                    column=len(line.rstrip()) + 1,
                    rule="trailing-whitespace",
                    severity="warning",
                    message="Trailing whitespace",
                    suggestion="Remove trailing whitespace"
                ))
                if fix:
                    fixed_line = line.rstrip()
            
            fixed_lines.append(fixed_line)
        
        # Check naming conventions
        naming_issues = self._check_python_naming(content, filepath)
        issues.extend(naming_issues)
        
        # Ensure final newline
        if fix and fixed_lines and fixed_lines[-1]:
            fixed_lines.append('')
        
        fixed_content = '\n'.join(fixed_lines)
        return issues, fixed_content
    
    def _check_python_naming(self, content: str, filepath: str) -> List[StyleIssue]:
        """Check Python naming conventions"""
        issues = []
        
        # Class names should be PascalCase
        class_pattern = r'class\s+([a-z_][a-zA-Z0-9_]*)\s*[\(:]'
        for match in re.finditer(class_pattern, content):
            class_name = match.group(1)
            if not class_name[0].isupper():
                line_num = content[:match.start()].count('\n') + 1
                issues.append(StyleIssue(
                    file=filepath,
                    line=line_num,
                    column=match.start() - content.rfind('\n', 0, match.start()),
                    rule="naming",
                    severity="warning",
                    message=f"Class '{class_name}' should be PascalCase",
                    suggestion=self._to_pascal_case(class_name)
                ))
        
        # Function names should be snake_case
        func_pattern = r'def\s+([A-Z][a-zA-Z0-9_]*)\s*\('
        for match in re.finditer(func_pattern, content):
            func_name = match.group(1)
            line_num = content[:match.start()].count('\n') + 1
            issues.append(StyleIssue(
                file=filepath,
                line=line_num,
                column=match.start() - content.rfind('\n', 0, match.start()),
                rule="naming",
                severity="warning",
                message=f"Function '{func_name}' should be snake_case",
                suggestion=self._to_snake_case(func_name)
            ))
        
        return issues
    
    def _check_javascript(self, content: str, filepath: str, fix: bool) -> Tuple[List[StyleIssue], str]:
        """Check JavaScript style issues"""
        issues = []
        lines = content.split('\n')
        fixed_lines = []
        
        for i, line in enumerate(lines, 1):
            fixed_line = line
            
            # Check indentation
            if line and line[0] in ' \t':
                if '\t' in line:
                    issues.append(StyleIssue(
                        file=filepath,
                        line=i,
                        column=1,
                        rule="indentation",
                        severity="error",
                        message="Tabs used for indentation",
                        suggestion="Use 2 spaces"
                    ))
                    if fix:
                        fixed_line = line.replace('\t', '  ')
            
            # Check trailing whitespace
            if line != line.rstrip():
                issues.append(StyleIssue(
                    file=filepath,
                    line=i,
                    column=len(line.rstrip()) + 1,
                    rule="trailing-whitespace",
                    severity="warning",
                    message="Trailing whitespace",
                    suggestion="Remove trailing whitespace"
                ))
                if fix:
                    fixed_line = line.rstrip()
            
            # Check semicolons
            trimmed = line.strip()
            if trimmed and not trimmed.startswith('//'):
                if re.match(r'^(const|let|var|return)\s+', trimmed):
                    if not trimmed.endswith((';', '{', '}')):
                        issues.append(StyleIssue(
                            file=filepath,
                            line=i,
                            column=len(line),
                            rule="semicolon",
                            severity="error",
                            message="Missing semicolon",
                            suggestion="Add semicolon"
                        ))
                        if fix:
                            fixed_line = line.rstrip() + ';'
            
            fixed_lines.append(fixed_line)
        
        # Check naming conventions
        naming_issues = self._check_javascript_naming(content, filepath)
        issues.extend(naming_issues)
        
        # Ensure final newline
        if fix and fixed_lines and fixed_lines[-1]:
            fixed_lines.append('')
        
        fixed_content = '\n'.join(fixed_lines)
        return issues, fixed_content
    
    def _check_javascript_naming(self, content: str, filepath: str) -> List[StyleIssue]:
        """Check JavaScript naming conventions"""
        issues = []
        
        # Variables should be camelCase
        var_pattern = r'(?:const|let|var)\s+([a-z]+_[a-z_]+)\s*='
        for match in re.finditer(var_pattern, content):
            var_name = match.group(1)
            line_num = content[:match.start()].count('\n') + 1
            issues.append(StyleIssue(
                file=filepath,
                line=line_num,
                column=match.start() - content.rfind('\n', 0, match.start()),
                rule="naming",
                severity="warning",
                message=f"Variable '{var_name}' should be camelCase",
                suggestion=self._to_camel_case(var_name)
            ))
        
        # Classes should be PascalCase
        class_pattern = r'class\s+([a-z][a-zA-Z0-9]*)'
        for match in re.finditer(class_pattern, content):
            class_name = match.group(1)
            line_num = content[:match.start()].count('\n') + 1
            issues.append(StyleIssue(
                file=filepath,
                line=line_num,
                column=match.start() - content.rfind('\n', 0, match.start()),
                rule="naming",
                severity="warning",
                message=f"Class '{class_name}' should be PascalCase",
                suggestion=self._to_pascal_case(class_name)
            ))
        
        return issues
    
    def _check_universal(self, content: str, filepath: str, fix: bool) -> Tuple[List[StyleIssue], str]:
        """Universal style checks for any file"""
        issues = []
        lines = content.split('\n')
        fixed_lines = []
        
        for i, line in enumerate(lines, 1):
            fixed_line = line
            
            # Check trailing whitespace
            if line != line.rstrip():
                issues.append(StyleIssue(
                    file=filepath,
                    line=i,
                    column=len(line.rstrip()) + 1,
                    rule="trailing-whitespace",
                    severity="warning",
                    message="Trailing whitespace",
                    suggestion="Remove trailing whitespace"
                ))
                if fix:
                    fixed_line = line.rstrip()
            
            fixed_lines.append(fixed_line)
        
        # Ensure final newline
        if fix and fixed_lines and fixed_lines[-1]:
            fixed_lines.append('')
        
        fixed_content = '\n'.join(fixed_lines)
        return issues, fixed_content
    
    def check_project(self, project_path: str, fix: bool = False, exclude: List[str] = None) -> List[StyleIssue]:
        """Check entire project for style issues"""
        path = Path(project_path)
        
        if not path.exists():
            print(f"Project path not found: {project_path}")
            return []
        
        exclude = exclude or self.config["ignore_paths"]
        
        for root, dirs, files in os.walk(path):
            # Remove excluded directories
            dirs[:] = [d for d in dirs if not any(ex in d for ex in exclude)]
            
            for file in files:
                filepath = os.path.join(root, file)
                # Check supported file types
                if any(filepath.endswith(ext) for ext in ['.py', '.js', '.jsx', '.ts', '.tsx', '.java', '.go', '.rb', '.php']):
                    self.check_file(filepath, fix)
        
        return self.issues
    
    def generate_report(self, format: str = "text") -> str:
        """Generate style report in specified format"""
        if format == "json":
            return self._generate_json_report()
        elif format == "markdown":
            return self._generate_markdown_report()
        else:
            return self._generate_text_report()
    
    def _generate_text_report(self) -> str:
        """Generate plain text report"""
        report = []
        report.append("=" * 60)
        report.append("CODE STYLE REPORT")
        report.append("=" * 60)
        report.append(f"Files checked: {self.files_checked}")
        report.append(f"Files fixed: {self.files_fixed}")
        report.append(f"Total issues: {len(self.issues)}")
        report.append(f"Issues fixed: {sum(1 for i in self.issues if i.fixed)}")
        report.append("")
        
        # Group issues by file
        by_file = {}
        for issue in self.issues:
            if issue.file not in by_file:
                by_file[issue.file] = []
            by_file[issue.file].append(issue)
        
        for filepath, file_issues in by_file.items():
            report.append(f"\n{filepath}")
            report.append("-" * len(filepath))
            
            for issue in file_issues:
                status = "✓ Fixed" if issue.fixed else "✗ Manual fix required"
                report.append(f"  Line {issue.line}: {issue.message} [{status}]")
                if issue.suggestion and not issue.fixed:
                    report.append(f"    Suggestion: {issue.suggestion}")
        
        return "\n".join(report)
    
    def _generate_json_report(self) -> str:
        """Generate JSON report"""
        report = {
            "summary": {
                "files_checked": self.files_checked,
                "files_fixed": self.files_fixed,
                "total_issues": len(self.issues),
                "issues_fixed": sum(1 for i in self.issues if i.fixed)
            },
            "issues": [asdict(issue) for issue in self.issues]
        }
        return json.dumps(report, indent=2)
    
    def _generate_markdown_report(self) -> str:
        """Generate Markdown report"""
        report = []
        report.append("# Code Style Report")
        report.append("")
        report.append("## Summary")
        report.append(f"- **Files Checked**: {self.files_checked}")
        report.append(f"- **Files Fixed**: {self.files_fixed}")
        report.append(f"- **Total Issues**: {len(self.issues)}")
        report.append(f"- **Issues Fixed**: {sum(1 for i in self.issues if i.fixed)}")
        report.append("")
        
        if self.issues:
            report.append("## Issues")
            
            # Group by severity
            by_severity = {}
            for issue in self.issues:
                if issue.severity not in by_severity:
                    by_severity[issue.severity] = []
                by_severity[issue.severity].append(issue)
            
            for severity in ["error", "warning", "info"]:
                if severity in by_severity:
                    emoji = {"error": "🔴", "warning": "🟡", "info": "🔵"}.get(severity, "")
                    report.append(f"\n### {emoji} {severity.upper()} ({len(by_severity[severity])})")
                    
                    for issue in by_severity[severity][:5]:  # Show first 5
                        status = "✅" if issue.fixed else "❌"
                        report.append(f"- {status} `{issue.file}:{issue.line}` - {issue.message}")
        
        return "\n".join(report)
    
    # Utility methods for name conversion
    def _to_camel_case(self, name: str) -> str:
        """Convert to camelCase"""
        parts = name.split('_')
        return parts[0].lower() + ''.join(p.capitalize() for p in parts[1:])
    
    def _to_pascal_case(self, name: str) -> str:
        """Convert to PascalCase"""
        parts = name.replace('_', ' ').split()
        return ''.join(p.capitalize() for p in parts)
    
    def _to_snake_case(self, name: str) -> str:
        """Convert to snake_case"""
        result = re.sub('([A-Z])', r'_\1', name)
        return result.lower().lstrip('_')

def main():
    parser = argparse.ArgumentParser(description="Code Style Enforcer")
    
    # Input options
    input_group = parser.add_mutually_exclusive_group(required=True)
    input_group.add_argument("--file", help="Check/fix single file")
    input_group.add_argument("--project", help="Check/fix entire project")
    
    # Action options
    parser.add_argument("--check", action="store_true", help="Check only (default)")
    parser.add_argument("--fix", action="store_true", help="Auto-fix issues")
    
    # Configuration
    parser.add_argument("--style", help="Style guide (pep8, black, airbnb, standard)")
    parser.add_argument("--config", help="Configuration file path")
    
    # Output options
    parser.add_argument("--report", help="Save report to file")
    parser.add_argument("--format", choices=["text", "json", "markdown"], default="text",
                       help="Report format")
    
    # Other options
    parser.add_argument("--exclude", help="Comma-separated paths to exclude")
    
    args = parser.parse_args()
    
    # Default to check mode
    if not args.check and not args.fix:
        args.check = True
    
    # Initialize enforcer
    enforcer = CodeStyleEnforcer(args.config)
    
    # Apply style guide overrides
    if args.style:
        if args.style == "black":
            enforcer.config["python"]["max_line_length"] = 88
        elif args.style == "pep8":
            enforcer.config["python"]["max_line_length"] = 79
        elif args.style == "airbnb":
            enforcer.config["javascript"]["indent_size"] = 2
        elif args.style == "standard":
            enforcer.config["javascript"]["semicolons"] = False
    
    # Process files
    fix = args.fix
    if args.file:
        enforcer.check_file(args.file, fix)
    elif args.project:
        exclude = args.exclude.split(",") if args.exclude else None
        enforcer.check_project(args.project, fix, exclude)
    
    # Generate report
    report = enforcer.generate_report(args.format)
    
    # Output report
    if args.report:
        with open(args.report, 'w') as f:
            f.write(report)
        print(f"Report saved to: {args.report}")
    else:
        print(report)
    
    # Exit with error if issues found and not fixed
    if enforcer.issues and not fix:
        sys.exit(1)

if __name__ == "__main__":
    main()