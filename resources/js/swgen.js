$(document).ready(function() {
	let urlpath = document.location.pathname;

	// show current links
	let link = $(`nav#content li a[href="${urlpath}"]`)
	link.addClass("li-focus");
	link.parents("details").attr("open",true);
	link[0].scrollIntoView();
})
