package listener

var html = `
<html>

<head>
    <link href="https://getgauge.io/assets/images/favicons/favicon.ico" rel="shortcut icon" type="image/ico">
    <link href="https://fonts.googleapis.com/css?family=Source+Code+Pro" rel="stylesheet">
    <title>Flash - {{ .Project }}</title>
</head>

<body>
    <div class="menu">
        <a href="http://getgauge.io" target="_blank" class="back"><img src="https://getgauge.io/assets/images/Gauge_logo.svg" draggable="false" /><span class="beta">BETA</span></a>
        <div class="toggleCheck"><input id="collapse" type="checkbox" onclick="toggleAll(this)" /> <label for="collapse">Collapse All</label></div>
        <div class="statsContainer">
            <table>
                <tr>
                    <td>Specification</td>
                    <td id="specPassedCount">0 ✔</td>
                    <td id="specFailedCount">0 ✘</td>
                    <td id="specProgressCount">0 ⚡</td>
                    <tr>
                        <td>Scenario&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</td>
                        <td id="scenarioPassedCount">0 ✔</td>
                        <td id="scenarioFailedCount">0 ✘</td>
                        <td id="scenarioProgressCount">0 ⚡</td>
            </table>
        </div>
    </div>
    <ul id="specs"></ul>
    <script type="text/javascript">
        const toggleAll = (element) => {
            document.querySelectorAll(".scenarios").forEach(function(e) {
                if (element.checked) e.hidden = true;
                else e.hidden = false;
            });
            document.querySelectorAll(".toggle").forEach(function(e) {
                if (element.checked) e.innerHTML = " +";
                else e.innerHTML = " -";
            });
        };

        const toggle = (id) => {
            var element = document.getElementById(id);
            var toggleElement = document.getElementById("toggle" + id);
            if (element.hidden) toggleElement.innerHTML = " -";
            else toggleElement.innerHTML = " +";
            element.hidden = !element.hidden;
        };

        const initWS = () => {
            const socket = new WebSocket("ws://{{ .URL }}");
            socket.onopen = () => console.log("Socket is open");
            socket.onmessage = function(e) {
                const data = JSON.parse(e.data);
                try {
                    if (data.Type == "end") document.getElementsByClassName("menu")[0].style.borderColor = data.Status == "fail" ? "#e73e48" : "#1cb596";
                    else if (data.Type == "spec") handleSpecEvent(data, data.Status);
                    else if (data.Type == "scenario") handleScenarioEvent(data, data.Status);
                    else handleStepEvent(data, data.Status);
                } finally {
                    updateCounts(counts);
                }
            };
            socket.onclose = () => console.log("Socket is closed");
            return socket;
        };

        const handleSpecEvent = (data, statusClass) => {
            const id = getID();
            if (specMap[data.FileName]) {
                document.getElementById(specMap[data.FileName].id).getElementsByTagName("li")[0].className = statusClass;
                counts.specProgress--;
                if (statusClass == "fail") counts.specFailed++;
                else counts.specPassed++;
                return;
            }
            const specs = document.getElementById("specs"),
                scnId = getID();
            const hideScenarios = document.getElementById("collapse").checked;
            specs.innerHTML += "<div class=\"spec\" id=\"" + id + "\"><li class=\"" + statusClass + "\"><span class=\"specName\">" + data.Name + "</span><span class=\"toggle\" id=\"toggle" + scnId + "\" onClick=\"toggle('" + scnId + "')\"> " + (hideScenarios ? "+" : "-") + "</span><ul" + (hideScenarios ? " hidden" : " ") + " class=\"scenarios\" id=\"" + scnId + "\"></ul></li></div>";
            specMap[data.FileName] = {
                "id": id,
                "scenarios": {},
            };
            counts.specProgress++;
        };

        const handleScenarioEvent = (data, statusClass) => {
            const id = getID();
            if (specMap[data.SpecFileName].scenarios[data.Name]) {
                const element = specMap[data.SpecFileName].scenarios[data.Name];
                document.getElementById(element.id).getElementsByTagName("li")[0].className = statusClass;
                counts.scenarioProgress--;
                if (statusClass == "fail") counts.scenarioFailed++;
                else counts.scenarioPassed++;
                return;
            }
            specMap[data.SpecFileName].scenarios[data.Name] = {
                "id": id,
                "steps": {},
            };
            const element = document.getElementById(specMap[data.SpecFileName].id).getElementsByTagName("ul")[0];
            element.insertAdjacentHTML("beforeend", "<div class=\"scenario\" id=\"" + id + "\"><li class=\"" + statusClass + "\"><span class=\"scenarioName\">" + data.Name + "</span><ul></ul></li></div>");
            counts.scenarioProgress++;
        };

        const handleStepEvent = (data, statusClass) => {
            const id = getID();
            if (specMap[data.SpecFileName].scenarios[data.ScenarioName].steps[data.Name]) {
                const element = specMap[data.SpecFileName].scenarios[data.ScenarioName].steps[data.Name];
                document.getElementById(element.id).getElementsByTagName("li")[0].className = statusClass;
                return;
            }
            specMap[data.SpecFileName].scenarios[data.ScenarioName].steps[data.Name] = {
                "id": id
            };
            const element = document.getElementById(specMap[data.SpecFileName].scenarios[data.ScenarioName].id).getElementsByTagName("ul")[0];
            element.insertAdjacentHTML("beforeend", "<div class=\"step\" id=\"" + id + "\"><li class=\"" + statusClass + "\"><span class=\"stepName\">" + data.Name + "</span><ul></ul></li></div>");
        };

        if (window.WebSocket === undefined) document.getElementById("specs").innerHTML = "Your browser does not support WebSockets";
        else initWS();
        const specMap = {};
        const counts = {
            "specFailed": 0,
            "specPassed": 0,
            "specProgress": 0,
            "scenarioFailed": 0,
            "scenarioPassed": 0,
            "scenarioProgress": 0
        };

        const getID = () => Math.random().toString(36).substring(7);

        const updateCounts = (counts) => {
            document.getElementById("specFailedCount").innerHTML = counts.specFailed + " ✘";
            document.getElementById("specPassedCount").innerHTML = counts.specPassed + " ✔";
            document.getElementById("specProgressCount").innerHTML = counts.specProgress + " ⚡";
            document.getElementById("scenarioFailedCount").innerHTML = counts.scenarioFailed + " ✘";
            document.getElementById("scenarioPassedCount").innerHTML = counts.scenarioPassed + " ✔";
            document.getElementById("scenarioProgressCount").innerHTML = counts.scenarioProgress + " ⚡";
        };
    </script>
    <style>
        body {
            font-family: 'Open Sans', sans-serif;
            margin: 0px;
            line-height: 20px;
            background: #f7f7f7;
        }

        .menu {
            width: 100%;
            height: 65px;
            background: #464545;
            z-index: 100;
            -webkit-touch-callout: none;
            -webkit-user-select: none;
            -moz-user-select: none;
            -ms-user-select: none;
            border-bottom: 4px solid #d0d032;
            position: fixed;
            margin-bottom: 1%;
            border-radius: 0px;
            text-align: center;
        }

        .back {
            position: absolute;
            width: 90px;
            height: 50px;
            top: 5px;
            left: 0px;
            color: #000000;
            line-height: 45px;
            font-size: 40px;
            padding-left: 10px;
            cursor: pointer;
            transition: .15s all;
            text-decoration: none;
        }

        .back img {
            position: absolute;
            top: 10px;
            left: 50px;
            height: 35px;
            margin-left: 15px;
        }

        .beta {
            position: relative;
            background: #f5c10e;
            color: #000;
            padding: 2px 6px;
            font-size: 12px;
            text-transform: uppercase;
            top: 6px;
            left: 177px;
            border-radius: 6px;
        }

        .back:active {
            background: rgba(0, 0, 0, 0.15);
        }

        .spec {
            font-size: 18px;
            margin-bottom: 0.8%;
            border-bottom: 1px solid #eaeaea;
            padding-bottom: 0.2%;
        }

        .spec:hover {
            background: #f1f1f1;
        }

        .scenario {
            font-size: 17px;
        }

        .step {
            font-size: 16px;
        }

        #specs {
            font-family: 'Source Code Pro', monospace;
            height: 85%;
            padding-top: 85px;
        }

        li {
            list-style: none;
        }

        ul {
            margin-left: -10px;
        }

        ul ul {
            margin-left: -20px;
        }

        ul ul ul {
            margin-left: -20px;
        }

        .toggle {
            color: gray;
            cursor: pointer;
        }

        .menu table td {
            color: #e2e2e2;
            font-weight: 100;
            padding-top: .5%;
            margin: 0px;
        }

        #specPassedCount,
        #scenarioPassedCount {
            color: #1cb596;
        }

        #specFailedCount,
        #scenarioFailedCount {
            color: #e73e48;
        }

        #specProgressCount,
        #scenarioProgressCount {
            color: #fcce3d;
        }

        .specName {
            padding-bottom: 0.1%;
            display: inline-block;
        }

        .scenarioName {
            padding-bottom: 0.05%;
            display: inline-block;
        }

        .stepName {
            padding-bottom: 0.05%;
            display: inline-block;
        }

        .toggleCheck {
            font-size: 17px;
            color: #e2e2e2;
            font-family: sans-serif;
            font-weight: 100;
            float: right;
            margin-right: 2%;
            padding-top: 25px;
        }

        .statsContainer {
            background: #464545;
            margin-top: 0.35%;
            margin-left: 41%;
        }

        table {
            text-align: center;
        }

        td {
            padding-right: 10px;
        }

        tr {
            line-height: 120%;
        }

        .fail {
            color: #d80a16;
        }

        .pass {
	    color: green;
        }

        .progress {
            color: black;
        }
    </style>
</body>

</html>
`
