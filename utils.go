package textwire

// errorPage returns HTML that's displayed when an error
// occurs while rendering a template
func errorPage() string {
	return `
	<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">	
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Something went wrong</title>
			</head>
		<body style="
			font-family: sans-serif;
			text-align: center;
			background: #101420;
			color: #eee;
		">
			<h1 style="
				margin-bottom: 0;
				color: #f4cb23;
				font-size: 4rem;
			">Oops!</h1>

			<p style="font-size: 1.3rem">
				Sorry! We’re having some trouble right now. Please check back shortly
			</p>

			<svg style="width: 170px; opacity: .05;" fill="none" stroke-width="1.5" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
				<path stroke-linecap="round" stroke-linejoin="round" d="M15.182 16.318A4.486 4.486 0 0 0 12.016 15a4.486 4.486 0 0 0-3.198 1.318M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Z"></path>
			</svg>
		</body>
	</html>
	`
}
