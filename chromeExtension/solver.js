var elements 	= document.getElementsByClassName('card');
if (elements == null || elements.length <= 0) {
	return;
}
var colors		= new Array(elements.length);
var selected	= new Array(elements.length);
for (var i = 0; i < elements.length; i++) {
	colors[i] = "";
	selected[i] = false;
}

var clickEvent = document.createEvent('MouseEvents');
clickEvent.initEvent('click', false, true);	

for (var i = 0; i < elements.length; i++) {
	if (!selected[i]) {
		elements[i].dispatchEvent(clickEvent);
		if (colors[i] == "") {
			colors[i] = elements[i].style.backgroundColor;
		}
	
		for (var j = i + 1; j < elements.length; j++) {
			if (!selected[i]) {
				elements[i].dispatchEvent(clickEvent);
			}
			if (!selected[j]) {
				if (colors[j] != "") {
					if (colors[i] == colors[j]) {
						elements[j].dispatchEvent(clickEvent);
						
						selected[i] = true;
						selected[j] = true;
						
						break;
					}
				} else {
					elements[j].dispatchEvent(clickEvent);
					colors[j] = elements[j].style.backgroundColor;
					if (colors[i] == colors[j]) {
					
						selected[i] = true;
						selected[j] = true;
						
						break;
					}
				}
			}		
		}
	}
}




/*
alert(elements.length);
return;

var element = document.getElementById('card0');
if (element == null) {
  alert('Card element is not found. Check element id.');
} else {
  var myevent = document.createEvent('MouseEvents');
  myevent.initEvent('click', false, true);
  element.dispatchEvent(myevent);
  alert('Card color is "' + element.style.backgroundColor + '".');
}
*/