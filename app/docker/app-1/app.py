from flask import Flask, request, redirect, url_for, render_template_string

app = Flask(__name__)
todos = []

TEMPLATE = '''
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Mini To-Do App</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.3/dist/css/bootstrap.min.css" rel="stylesheet">
</head>
<body class="bg-light">
<div class="container py-5">
    <h1 class="mb-4 text-center">📝 To-Do List</h1>
    <form method="POST" action="/" class="d-flex mb-3">
        <input type="text" name="item" class="form-control me-2" placeholder="Add new task..." required>
        <button class="btn btn-primary">Add</button>
    </form>
    <ul class="list-group">
        {% for todo in todos %}
            <li class="list-group-item d-flex justify-content-between align-items-center">
                {{ todo }}
                <a href="{{ url_for('delete', index=loop.index0) }}" class="btn btn-sm btn-danger">✖</a>
            </li>
        {% endfor %}
    </ul>
</div>
</body>
</html>
'''

@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        todos.append(request.form['item'])
        return redirect(url_for('index'))
    return render_template_string(TEMPLATE, todos=todos)

@app.route('/delete/<int:index>')
def delete(index):
    if 0 <= index < len(todos):
        todos.pop(index)
    return redirect(url_for('index'))

if __name__ == '__main__':
    app.run(debug=True)
