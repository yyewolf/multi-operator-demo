from flask import Flask, request, redirect, url_for, render_template_string
import random

app = Flask(__name__)
metrics = {
    "Visitors": random.randint(100, 1000),
    "Sales": random.randint(10, 100),
    "Errors": random.randint(0, 10)
}

TEMPLATE = '''
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Mini Dashboard</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
</head>
<body>
<section class="section">
  <div class="container">
    <h1 class="title has-text-centered">📈 Mini Dashboard</h1>
    <div class="columns is-multiline mt-5">
      {% for label, value in metrics.items() %}
        <div class="column is-one-third">
          <div class="box has-text-centered">
            <p class="heading">{{ label }}</p>
            <p class="title">{{ value }}</p>
          </div>
        </div>
      {% endfor %}
    </div>
    <form method="POST" class="has-text-centered mt-5">
      <button class="button is-info is-light">🔄 Refresh Data</button>
    </form>
  </div>
</section>
</body>
</html>
'''

@app.route('/', methods=['GET', 'POST'])
def dashboard():
    if request.method == 'POST':
        for key in metrics:
            metrics[key] = random.randint(0, 1000)
        return redirect(url_for('dashboard'))
    return render_template_string(TEMPLATE, metrics=metrics)

if __name__ == '__main__':
    app.run(debug=True)
