import re
import zipfile
from jinja2 import Template

# Define your keywords and regular expression patterns
keywords = {
    'email': r'srv.*AP\d+[A-Za-z]*',
    'password': r'.*?(?:secret|password).*?"(.*?)"',
    # Add more keywords and patterns as needed
}

def search_in_zip(zip_filename, keyword_patterns):
    results = {keyword: [] for keyword in keyword_patterns.keys()}
    
    with zipfile.ZipFile(zip_filename, 'r') as zip_ref:
        for file_info in zip_ref.infolist():
            with zip_ref.open(file_info) as file:
                content = file.read().decode('utf-8')
                for keyword, pattern in keyword_patterns.items():
                    matches = re.findall(pattern, content, re.IGNORECASE)
                    if matches:
                        results[keyword].extend([(match, file_info.filename) for match in matches])
    
    return results

def generate_html_report(results):
    template_str = """
    <!DOCTYPE html>
    <html>
    <head>
        <title>Keyword Search Report</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                margin: 0;
                padding: 0;
                display: flex;
                justify-content: flex-start;
            }
            .container {
                max-width: 800px;
                padding: 20px;
            }
            .keyword-group {
                border: 1px solid #ddd;
                padding: 10px;
                margin: 10px 0;
                background-color: #f9f9f9;
            }
            .keyword-group h2 {
                margin: 0;
                font-size: 1.2em;
            }
            .results-list {
                margin-top: 10px;
                padding-left: 20px;
            }
            .file-info {
                font-size: 0.8em;
                color: #666;
            }
            .keyword-buttons {
                display: flex;
                justify-content: flex-start;
                margin-top: 10px;
            }
            .keyword-button {
                margin-right: 10px;
                background-color: #ccc;
                padding: 5px 10px;
                cursor: pointer;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="keyword-buttons">
                {% for keyword in results.keys() %}
                <button class="keyword-button" onclick="toggleResults('{{ keyword }}')">{{ keyword | title }}</button>
                {% endfor %}
            </div>
            {% for keyword, matches in results.items() %}
            <div id="{{ keyword }}-results" class="results-list" style="display: none;">
                {% if matches %}
                {% for match, filepath in matches %}
                <p>{{ match }} <span class="file-info">({{ filepath }})</span></p>
                {% endfor %}
                {% else %}
                <p>No {{ keyword }} matches found.</p>
                {% endif %}
            </div>
            {% endfor %}
        </div>
        <script>
            function toggleResults(keyword) {
                {% for keyword in results.keys() %}
                var {{ keyword }}ResultsList = document.getElementById('{{ keyword }}-results');
                {% endfor %}
                
                {% for keyword in results.keys() %}
                if (keyword === '{{ keyword }}') {
                    {{ keyword }}ResultsList.style.display = ({{ keyword }}ResultsList.style.display === "none") ? "block" : "none";
                } else {
                    {{ keyword }}ResultsList.style.display = "none";
                }
                {% endfor %}
            }
        </script>
    </body>
    </html>
    """
    
    template = Template(template_str)
    report_html = template.render(results=results)
    return report_html

if __name__ == "__main__":
    zip_filename = "w.zip"  # Replace with your zip file's name
    keyword_results = search_in_zip(zip_filename, keywords)
    report_html = generate_html_report(keyword_results)
    
    with open("report.html", "w", encoding="utf-8") as report_file:
        report_file.write(report_html)
    
    print("Report generated successfully as 'report.html'")
