<script>
function updateEssayList(str)
{
var xmlhttp;    
if (str=="")
  {
  document.getElementById("essaylist").innerHTML="";
  return;
  }
if (window.XMLHttpRequest)
  {// code for IE7+, Firefox, Chrome, Opera, Safari
  xmlhttp=new XMLHttpRequest();
  }
else
  {// code for IE6, IE5
  xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
  }
xmlhttp.onreadystatechange=function()
  {    
  if (xmlhttp.readyState==4 && xmlhttp.status==200)
    {
    document.getElementById("essaylistTest").innerHTML=xmlhttp.responseText+str;
    }
  }  
  // try this for more times!!!! success!!!
xmlhttp.open("GET","{{url "Essay.PageList" ""}}"+str,true);
xmlhttp.send();
}
</script> 