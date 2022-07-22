<h1 id="about-this-software">About this software</h1>
<h2 id="name-version">${name} ${version}</h2>
<p>${if(author)} ### Authors</p>
<p>${for(author)} - ${it.givenName} ${it.familyName} ${endfor}
${endif}</p>
<p>${if(description)} ${description} ${endif}</p>
<p>${if(license)}- License: <span
class="math inline"><em>l</em><em>i</em><em>c</em><em>e</em><em>n</em><em>s</em><em>e</em></span>{endif}
${if(codeRepository)}- GitHub: <span
class="math inline"><em>c</em><em>o</em><em>d</em><em>e</em><em>R</em><em>e</em><em>p</em><em>o</em><em>s</em><em>i</em><em>t</em><em>o</em><em>r</em><em>y</em></span>{endif}
${if(issueTracker)}- Issues: <span
class="math inline"><em>i</em><em>s</em><em>s</em><em>u</em><em>e</em><em>T</em><em>r</em><em>a</em><em>c</em><em>k</em><em>e</em><em>r</em></span>{endif}</p>
<p>${if(programmingLanguage)} ### Programming languages</p>
<p>${for(programmingLanguage)} - ${programmingLanguage} ${endfor}
${endif}</p>
<p>${if(operatingSystem)} ### Operating Systems</p>
<p>${for(operatingSystem)} - ${operatingSystem} ${endfor} ${endif}</p>
<p>${if(softwareRequirements)} ### Software Requiremets</p>
<p>${for(softwareRequirements)} - ${softwareRequirements} ${endfor}
<span
class="math inline"><em>e</em><em>n</em><em>d</em><em>i</em><em>f</em></span></p>

