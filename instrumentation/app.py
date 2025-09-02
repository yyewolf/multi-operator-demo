from flask import Flask, redirect, url_for, render_template_string
import subprocess
import json
import time

app = Flask(__name__)

TEMPLATE = '''
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>K8s Presentation Control</title>
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bulma@0.9.4/css/bulma.min.css">
</head>
<body>
<section class="section">
    <div class="container">
        <h1 class="title has-text-centered mb-6">🎛️ Presentation Control Panel</h1>
        <div class="columns is-centered">
            <div class="column is-4">
                <h2 class="subtitle has-text-centered">Agent Operations</h2>
                <form action="/start" method="POST">
                    <button class="button is-success is-fullwidth mb-2">▶️ Start App</button>
                </form>
                <form action="/stop" method="POST">
                    <button class="button is-danger is-fullwidth mb-2">⏹ Stop App</button>
                </form>
                <form action="/restart" method="POST">
                    <button class="button is-warning is-fullwidth">🔁 Restart App</button>
                </form>
            </div>
            <div class="column is-4">
                <h2 class="subtitle has-text-centered">Resource Modifications</h2>
                <form action="/state/1" method="POST">
                    <button class="button is-link is-light is-fullwidth mb-2">Set Command to App 1</button>
                </form>
                <form action="/state/2" method="POST">
                    <button class="button is-link is-light is-fullwidth mb-2">Set Command to App 2</button>
                </form>
                <form action="/state/3" method="POST">
                    <button class="button is-link is-light is-fullwidth mb-2">Set Route to Maintenance</button>
                </form>
                <form action="/reload" method="POST">
                    <button class="button is-link is-light is-fullwidth">Force Reload App Config</button>
                </form>
            </div>
        </div>
    </div>
</section>
</body>
</html>
'''

def run(cmd):
    try:
        result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
        print(f"[CMD] {cmd}\n{result.stdout or result.stderr}")
    except Exception as e:
        print(f"Error: {e}")

def update_supervisor_in_all_pods():
    print("🔧 Updating supervisor in all pods...")
    try:
        pod_data = subprocess.run(
            ["kubectl", "get", "pods", "-n", "default", "-o", "json"],
            capture_output=True, text=True, check=True
        )
        pods = json.loads(pod_data.stdout)["items"]
        for pod in pods:
            pod_name = pod["metadata"]["name"]
            print(f"📦 Updating pod: {pod_name}")
    
            # Step 1: Annotate to trigger CM refresh
            timestamp = str(int(time.time()))
            annotate_cmd = (
                f"kubectl annotate pod -n default {pod_name} "
                f"presentation-refresh-timestamp={timestamp} --overwrite"
            )
            run(annotate_cmd)

            for cmd in [
                "supervisorctl -c /etc/supervisor.d/supervisord.conf reread",
                "supervisorctl -c /etc/supervisor.d/supervisord.conf update"
            ]:
                full_cmd = f"kubectl exec -n default {pod_name} -- {cmd}"
                run(full_cmd)
    except Exception as e:
        print(f"Failed to update supervisor: {e}")

@app.route('/')
def index():
    return render_template_string(TEMPLATE)

@app.route('/start', methods=['POST'])
def start_app():
    run("kubectl get --raw=/apis/agent.app.multi.ch/v1/namespaces/default/apps/app-sample/start")
    return redirect(url_for('index'))

@app.route('/stop', methods=['POST'])
def stop_app():
    run("kubectl get --raw=/apis/agent.app.multi.ch/v1/namespaces/default/apps/app-sample/stop")
    return redirect(url_for('index'))

@app.route('/restart', methods=['POST'])
def restart_app():
    run("kubectl get --raw=/apis/agent.app.multi.ch/v1/namespaces/default/apps/app-sample/restart")
    return redirect(url_for('index'))

@app.route('/state/<int:num>', methods=['POST'])
def apply_state(num):
    run(f"kubectl apply -k states/{num}")
    # sleep cause reconciliation takes time on purpose for the demonstration
    time.sleep(20)
    update_supervisor_in_all_pods()
    return redirect(url_for('index'))

    
@app.route('/reload', methods=['POST'])
def reload_config():
    update_supervisor_in_all_pods()
    return redirect(url_for('index'))

if __name__ == '__main__':
    app.run(debug=True)
