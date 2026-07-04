package htmltree

import (
	"bytes"
	"encoding/json"
	"html/template"

	"butter/pkg/ast"
	"butter/pkg/output"
)

func init() {
	output.Register(htmlTreeExt{})
}

type htmlTreeExt struct{}

func (htmlTreeExt) Name() string          { return "htmltree" }
func (htmlTreeExt) FileExtension() string { return ".html" }

func (htmlTreeExt) Serialize(spec *ast.AppSpec) ([]byte, error) {
	specJSON, err := json.Marshal(spec)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]interface{}{
		"Spec": template.JS(specJSON),
		"App":  spec.App,
	}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var tmpl = template.Must(template.New("htmltree").Parse(page))

const page = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Butter Spec — {{.App}}</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{background:#121212;background-image:linear-gradient(rgba(255,255,255,.03)1px,transparent 1px),linear-gradient(90deg,rgba(255,255,255,.03)1px,transparent 1px);background-size:20px 20px;font-family:'Courier New',Courier,monospace;color:#b0b0b0;overflow:hidden;height:100vh;width:100vw;position:fixed;-webkit-user-select:none;user-select:none;-webkit-touch-callout:none;touch-action:manipulation}
#svg-canvas{position:fixed;top:0;left:0;width:100vw;height:100vh;pointer-events:none;z-index:1}
.tree-panel{position:fixed;inset:0;overflow:hidden;cursor:grab;z-index:2}
.tree-panel.dragging{cursor:grabbing}
#viewport{transform-origin:0 0;min-width:100vw;min-height:100vh;display:flex;align-items:center;justify-content:center;padding:60px}
.tree-container{display:flex;align-items:flex-start;gap:100px}
.column{display:flex;flex-direction:column;gap:30px;justify-content:center;transition:opacity .3s ease}
.node-card{background:#1e1e1e;border:1px solid #333;border-radius:6px;width:240px;box-shadow:0 4px 15px rgba(0,0,0,.5);font-size:13px;overflow:hidden;transition:all .2s ease;flex-shrink:0}
.node-card:hover{border-color:#555}

.node-row{padding:8px 12px;border-bottom:1px solid #2a2a2a;display:flex;align-items:flex-start;gap:6px;word-break:break-word;overflow-wrap:break-word;white-space:normal;line-height:1.4}
.node-row:last-child{border-bottom:none}
.key{color:#D4A24C;margin-right:2px;word-break:break-word;overflow-wrap:break-word}
.val-string{color:#7ABA7A}
.val-number{color:#d19a66}
.val-bool{color:#8A7ADA}
.interactive-row{cursor:pointer;user-select:none}
.interactive-row:hover{background:#252525}
.toggle-icon{display:inline-block;width:10px;color:#64748B;flex-shrink:0;text-align:center;margin-top:2px}
.node-title{font-size:14px;font-weight:700;color:#5A9ABA;padding:8px 12px;border-bottom:1px solid #2a2a2a;display:flex;align-items:center;gap:8px;word-break:break-word}
.node-title .badge{font-size:10px;font-weight:400;padding:1px 6px;border-radius:3px;background:#333;color:#888;margin-left:auto;flex-shrink:0}
.desc-row{padding-left:20px;padding-bottom:8px;border-bottom:1px solid #2a2a2a;font-size:12px}
.desc-row .val-string{color:#b0b0b0;display:block}
.act-num{color:#64748B;min-width:18px;flex-shrink:0}
.param-key-row{padding:4px 12px 2px 16px;border-left:2px solid #777;border-bottom:none;background:#191919}
.param-key-row .key{color:#AA7A9A;font-size:12px}
.param-val-row{padding:2px 12px 4px 16px;border-left:2px solid #777;border-bottom:1px solid #222;background:#232323}
.param-val-row .val-string{color:#b0b0b0;font-size:11px}
.cond-key-row{padding:4px 12px 2px 20px;border-left:2px solid #555;border-bottom:none;background:#191919}
.cond-key-row .key{color:#8A7ADA;font-size:11px}
.cond-val-row{padding:2px 12px 4px 20px;border-left:2px solid #555;border-bottom:1px solid #222;background:#232323}
.cond-val-row .val-string{color:#7ABA7A;font-size:11px}
.enforce-key-row{padding:4px 12px 2px 20px;border-left:2px solid #555;border-bottom:none;background:#191919}
.enforce-key-row .key{color:#e11d48;font-size:11px}
.enforce-val-row{padding:2px 12px 8px 20px;border-left:2px solid #555;border-bottom:1px solid #222;background:#232323}
.enforce-val-row .val-string{color:#e6e6e6;font-size:11px}
.connection-path{fill:none;stroke:#555;stroke-width:1.5}
.connection-label{font-size:10px;fill:#666;font-family:'Courier New',Courier,monospace;user-select:none;pointer-events:none}
</style>
</head>
<body>
<svg id="svg-canvas"></svg>
<div class="tree-panel" id="tree-wrap">
  <div id="viewport">
    <div class="tree-container" id="tree-container">
      <div class="column" id="col-1"></div>
      <div class="column" id="col-2"></div>
      <div class="column" id="col-3"></div>
    </div>
  </div>
</div>
<script id="spec-data" type="application/json">{{.Spec}}</script>
<script>
(function(){
var spec = JSON.parse(document.getElementById('spec-data').textContent);
var state = { collapsed: {} };
state.collapsed['app'] = true;

function esc(s){var d=document.createElement('div');d.textContent=s;return d.innerHTML}

function buildTree(spec){
  var tree = {
    id: 'app',
    title: spec.app,
    details: [],
    children: []
  };
  if(spec.description) tree.details.push({k:'description',v:spec.description,type:'desc'});
  if(spec.version) tree.details.push({k:'version',v:spec.version});

  (spec.features||[]).forEach(function(f,fi){
    var feat = {
      id: 'f-'+fi,
      title: f.name,
      details: [],
      params: [],
      actions: [],
      children: []
    };
    if(f.description) feat.details.push({k:'description',v:f.description,type:'desc'});
    if(f.version) feat.details.push({k:'version',v:f.version});
     (f.params||[]).forEach(function(p){
      var parts = [p.type];
      if(p.required) parts.push('required');
      if(p.default !== undefined && p.default !== null && p.default !== '') parts.push('default: '+p.default);
      if(p.validate && p.validate.length) parts.push('validate: '+p.validate.join(', '));
      if(p.length) parts.push('length: '+p.length);
      feat.params.push({k:p.name,v:parts.join(', ')});
    });
    (f.actions||[]).forEach(function(a,ai){
      var act = {number:ai+1,statement:a.statement,conditions:[],enforces:[]};
      if(a.condition) act.conditions.push({type:a.condition.type,expr:a.condition.expression});
      if(a.enforce) a.enforce.forEach(function(e){act.enforces.push(e);});
      feat.actions.push(act);
    });
    feat.details.push({k:'params',v:(f.params||[]).length+' items',sub:'params-'+fi,type:'sub'});
    feat.details.push({k:'actions',v:(f.actions||[]).length+' items',sub:'actions-'+fi,type:'sub'});
    state.collapsed['params-'+fi] = true;
    state.collapsed['actions-'+fi] = true;
    tree.children.push(feat);
  });
  return tree;
}

function renderTree(tree){
  var col1 = document.getElementById('col-1');
  var col2 = document.getElementById('col-2');
  var col3 = document.getElementById('col-3');
  col1.innerHTML = '';
  col2.innerHTML = '';
  col3.innerHTML = '';

  var rootIcon = state.collapsed['app'] ? '+' : '&minus;';
  var rootHTML = '<div class="node-card" id="root-node">';
  rootHTML += '<div class="node-title">'+esc(tree.title)+'</div>';
  tree.details.forEach(function(d){
    if(d.type === 'desc'){
      rootHTML += '<div class="node-row"><span class="key">'+esc(d.k)+':</span></div>';
      rootHTML += '<div class="node-row desc-row"><span class="val-string">'+esc(d.v)+'</span></div>';
    } else {
      rootHTML += '<div class="node-row"><span class="key">'+esc(d.k)+':</span><span class="val-string">'+esc(d.v)+'</span></div>';
    }
  });
  rootHTML += '<div class="node-row interactive-row" onclick="window.__toggleRoot()">';
  rootHTML += '<span class="toggle-icon">'+rootIcon+'</span>';
  rootHTML += '<span class="key">features:</span>';
  rootHTML += '<span class="val-string">['+tree.children.length+' items]</span>';
  rootHTML += '</div></div>';
  col1.innerHTML = rootHTML;

  tree.children.forEach(function(feat){
    var featHTML = '<div class="node-card" id="card-'+feat.id+'">';
    featHTML += '<div class="node-title">'+esc(feat.title)+'</div>';
    feat.details.forEach(function(d){
      if(d.type === 'sub'){
        var subCollapsed = state.collapsed[d.sub];
        var subIcon = subCollapsed ? '+' : '&minus;';
        featHTML += '<div class="node-row interactive-row" data-sub="'+d.sub+'" onclick="window.__toggleSub(\''+d.sub+'\')">';
        featHTML += '<span class="toggle-icon">'+subIcon+'</span>';
        featHTML += '<span class="key">'+esc(d.k)+':</span>';
        featHTML += '<span class="val-string">'+esc(d.v)+'</span>';
        featHTML += '</div>';
      } else if(d.type === 'desc'){
        featHTML += '<div class="node-row"><span class="key">'+esc(d.k)+':</span></div>';
        featHTML += '<div class="node-row desc-row"><span class="val-string">'+esc(d.v)+'</span></div>';
      } else {
        featHTML += '<div class="node-row"><span class="key">'+esc(d.k)+':</span><span class="val-string">'+esc(d.v)+'</span></div>';
      }
    });
    featHTML += '</div>';
    col2.innerHTML += featHTML;

    var subHTML = '';
    if(feat.params.length){
      var pcollapsed = state.collapsed['params-'+feat.id.replace('f-','')];
      subHTML += '<div class="node-card" id="sub-params-'+feat.id.replace('f-','')+'"'+(pcollapsed?' style="display:none"':'')+'>';
      subHTML += '<div class="node-title">Params <span class="badge">'+feat.params.length+'</span></div>';
        feat.params.forEach(function(p){
          subHTML += '<div class="node-row param-key-row"><span class="key" style="color:#AA7A9A">'+esc(p.k)+':</span></div>';
          subHTML += '<div class="node-row param-val-row"><span class="val-string">'+esc(p.v)+'</span></div>';
        });
      subHTML += '</div>';
    }
    if(feat.actions.length){
      var acollapsed = state.collapsed['actions-'+feat.id.replace('f-','')];
      subHTML += '<div class="node-card" id="sub-actions-'+feat.id.replace('f-','')+'"'+(acollapsed?' style="display:none"':'')+'>';
      subHTML += '<div class="node-title">Actions <span class="badge">'+feat.actions.length+'</span></div>';
      feat.actions.forEach(function(a){
        subHTML += '<div class="node-row"><span class="act-num">'+a.number+'.</span><span class="val-string">'+esc(a.statement)+'</span></div>';
        a.conditions.forEach(function(c){
          subHTML += '<div class="node-row cond-key-row"><span class="key">'+esc(c.type)+':</span></div>';
          subHTML += '<div class="node-row cond-val-row"><span class="val-string">"'+esc(c.expr)+'"</span></div>';
        });
        a.enforces.forEach(function(e){
          subHTML += '<div class="node-row enforce-key-row"><span class="key">enforce:</span></div>';
          subHTML += '<div class="node-row enforce-val-row"><span class="val-string">'+esc(e)+'</span></div>';
        });
      });
      subHTML += '</div>';
    }
    col3.innerHTML += subHTML;
  });
}

function applySubVisibility(){
  var ftIdx=0;
  while(true){
    var featCard=document.getElementById('card-f-'+ftIdx);
    if(!featCard) break;
    var fi=ftIdx;
    var ps=document.getElementById('sub-params-'+fi);
    if(ps) ps.style.display=state.collapsed['params-'+fi]?'none':'block';
    var asub=document.getElementById('sub-actions-'+fi);
    if(asub) asub.style.display=state.collapsed['actions-'+fi]?'none':'block';
    ftIdx++;
  }
  var col3=document.getElementById('col-3');
  var anyVisible=false;
  document.querySelectorAll('#col-3 .node-card').forEach(function(c){if(c.style.display!=='none')anyVisible=true;});
  col3.style.opacity=anyVisible?'1':'0';
}

function drawConnections(){
  var svg = document.getElementById('svg-canvas');
  svg.innerHTML = '';
  svg.setAttribute('width', window.innerWidth);
  svg.setAttribute('height', window.innerHeight);

  var col2 = document.getElementById('col-2');
  var col3 = document.getElementById('col-3');
  if(state.collapsed['app']){
    col2.style.display = 'none';
    col3.style.display = 'none';
    return;
  }
  col2.style.display = 'flex';
  col3.style.display = 'flex';
  applySubVisibility();

  var rootCard = document.getElementById('root-node');
  if(!rootCard) return;

  var rl = rootCard.getBoundingClientRect();
  var ftIdx = 0;
  while(true){
    var featCard = document.getElementById('card-f-'+ftIdx);
    if(!featCard) break;
    var fl = featCard.getBoundingClientRect();
    drawBezier(rl.right, rl.top+rl.height/2, fl.left, fl.top+fl.height/2, svg);
    var fi = ftIdx;
    var ps = document.getElementById('sub-params-'+fi);
    if(ps&&!state.collapsed['params-'+fi]){
      var tr = featCard.querySelector('[data-sub="params-'+fi+'"]');
      if(tr){var tl=tr.getBoundingClientRect();var sl=ps.getBoundingClientRect();
        drawBezier(tl.right, tl.top+tl.height/2, sl.left, sl.top+sl.height/2, svg);
        drawLabel((tl.right+sl.left)/2,(tl.top+tl.height/2+sl.top+sl.height/2)/2-6,'params',svg);}
    }
    var asub = document.getElementById('sub-actions-'+fi);
    if(asub&&!state.collapsed['actions-'+fi]){
      var tr = featCard.querySelector('[data-sub="actions-'+fi+'"]');
      if(tr){var tl=tr.getBoundingClientRect();var sl=asub.getBoundingClientRect();
        drawBezier(tl.right, tl.top+tl.height/2, sl.left, sl.top+sl.height/2, svg);
        drawLabel((tl.right+sl.left)/2,(tl.top+tl.height/2+sl.top+sl.height/2)/2-6,'actions',svg);}
    }
    ftIdx++;
  }
}

function drawBezier(x1,y1,x2,y2,svg){
  var off = Math.abs(x2-x1)*0.5;
  var d = 'M '+x1+' '+y1+' C '+(x1+off)+' '+y1+', '+(x2-off)+' '+y2+', '+x2+' '+y2;
  var p = document.createElementNS('http://www.w3.org/2000/svg','path');
  p.setAttribute('d',d);p.setAttribute('class','connection-path');svg.appendChild(p);
}

function drawLabel(x,y,text,svg){
  var t = document.createElementNS('http://www.w3.org/2000/svg','text');
  t.setAttribute('x',x);t.setAttribute('y',y);t.setAttribute('text-anchor','middle');
  t.setAttribute('class','connection-label');t.textContent=text;svg.appendChild(t);
}

function applyColState(){
  var c2 = document.getElementById('col-2'), c3 = document.getElementById('col-3');
  c2.style.display = state.collapsed['app']?'none':'flex';
  c3.style.display = state.collapsed['app']?'none':'flex';
}

window.__toggleRoot = function(){
  state.collapsed['app'] = !state.collapsed['app'];
  var icon = document.querySelector('#root-node .toggle-icon');
  if(icon) icon.innerHTML = state.collapsed['app']?'+':'&minus;';
  applyColState(); drawConnections();
  setTimeout(zoomToFit, 0);
};

window.__toggleSub = function(key){
  state.collapsed[key] = !state.collapsed[key];
  var tr = document.querySelector('[data-sub="'+key+'"]');
  if(tr){var icon=tr.querySelector('.toggle-icon');if(icon)icon.innerHTML=state.collapsed[key]?'+':'&minus;';}
  drawConnections();
};

var zoom=1,panX=0,panY=0,dragging=false,dsx=0,dsy=0;
var vp=document.getElementById('viewport');
var wrap=document.getElementById('tree-wrap');

function updateTransform(){
  vp.style.transform='scale('+zoom+') translate('+(panX/zoom)+'px,'+(panY/zoom)+'px)';
  drawConnections();
}

wrap.addEventListener('wheel',function(e){
  e.preventDefault();
  var r=wrap.getBoundingClientRect(),mx=e.clientX-r.left,my=e.clientY-r.top,oz=zoom;
  zoom=Math.max(0.1,Math.min(5,zoom+(e.deltaY>0?-0.12:0.12)));
  var s=zoom/oz;panX=mx-s*(mx-panX);panY=my-s*(my-panY);updateTransform();
},{passive:false});

wrap.addEventListener('mousedown',function(e){
  if(e.button!==0)return;
  dragging=true;dsx=e.clientX-panX;dsy=e.clientY-panY;
  wrap.classList.add('dragging');
});

window.addEventListener('mousemove',function(e){
  if(!dragging)return;
  panX=e.clientX-dsx;panY=e.clientY-dsy;updateTransform();
});

window.addEventListener('mouseup',function(){dragging=false;wrap.classList.remove('dragging');});

var touching=false,tdx=0,tdy=0,tpx=0,tpy=0,pinchDist=0,pinchZoom=0;
wrap.addEventListener('touchstart',function(e){
  if(e.touches.length===1){
    touching=true;tdx=e.touches[0].clientX-panX;tdy=e.touches[0].clientY-panY;
    tpx=e.touches[0].clientX;tpy=e.touches[0].clientY;return;
  }
  if(e.touches.length===2){
    pinchDist=Math.hypot(e.touches[0].clientX-e.touches[1].clientX,e.touches[0].clientY-e.touches[1].clientY);
    pinchZoom=zoom;
  }
},{passive:true});

wrap.addEventListener('touchmove',function(e){
  if(e.touches.length===1&&touching){
    panX=e.touches[0].clientX-tdx;panY=e.touches[0].clientY-tdy;
    var s=e.touches[0].clientX-tpx;
    if(Math.abs(s)>5||Math.abs(e.touches[0].clientY-tpy)>5)e.preventDefault();
    tpx=e.touches[0].clientX;tpy=e.touches[0].clientY;
    updateTransform();return;
  }
  if(e.touches.length===2&&pinchDist>0){
    e.preventDefault();
    var d=Math.hypot(e.touches[0].clientX-e.touches[1].clientX,e.touches[0].clientY-e.touches[1].clientY);
    zoom=Math.max(0.1,Math.min(5,pinchZoom*(d/pinchDist)));
    var cx=(e.touches[0].clientX+e.touches[1].clientX)/2,cy=(e.touches[0].clientY+e.touches[1].clientY)/2;
    var r=wrap.getBoundingClientRect(),mx=cx-r.left,my=cy-r.top,oz=pinchZoom;
    var s=zoom/oz;panX=mx-s*(mx-panX);panY=my-s*(my-panY);
    updateTransform();
  }
},{passive:false});

wrap.addEventListener('touchend',function(){touching=false;pinchDist=0;wrap.classList.remove('dragging');});
wrap.addEventListener('touchcancel',function(){touching=false;pinchDist=0;wrap.classList.remove('dragging');});

var lastTap=0;
wrap.addEventListener('click',function(e){
  if(e.target.closest('.interactive-row,.node-title')) return;
  var now=Date.now();
  if(now-lastTap<300){
    var r=wrap.getBoundingClientRect(),mx=e.clientX-r.left,my=e.clientY-r.top;
    if(zoom<1.5){zoom=2;panX=mx-2*(mx-panX);panY=my-2*(my-panY);}
    else{zoom=1;panX=0;panY=0;}
    updateTransform();
  }
  lastTap=now;
});

var tree=buildTree(spec);
renderTree(tree);
applyColState();
drawConnections();

function zoomToFit(){
  var cards=document.querySelectorAll('.node-card');
  if(!cards.length){updateTransform();return;}
  var minL=Infinity,maxR=-Infinity,minT=Infinity,maxB=-Infinity;
  cards.forEach(function(c){
    if(c.offsetParent===null) return;
    var r=c.getBoundingClientRect();
    if(r.left<minL)minL=r.left;if(r.right>maxR)maxR=r.right;
    if(r.top<minT)minT=r.top;if(r.bottom>maxB)maxB=r.bottom;
  });
  if(minL===Infinity){updateTransform();return;}
  var wr=wrap.getBoundingClientRect(),pad=60;
  var w=maxR-minL+pad*2,h=maxB-minT+pad*2;
  var sx=wr.width/w,sy=wr.height/h;
  zoom=Math.min(sx,sy,0.95);zoom=Math.max(0.1,Math.min(5,zoom));
  panX=0;panY=0;updateTransform();
}

setTimeout(zoomToFit, 50);
window.addEventListener('resize',function(){setTimeout(zoomToFit,50);});
})();
</script>
</body>
</html>`
