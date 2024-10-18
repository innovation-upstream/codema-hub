package pagetemplate

const BountiesPageTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Bounties - CodemaHub</title>
    <link href="/static/output.css" rel="stylesheet">
		<script>
			// Check for dark mode preference
			if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
				document.documentElement.classList.add('dark')
			} else {
				document.documentElement.classList.remove('dark')
			}
		</script>
	</head>
	<body class="bg-gray-100 dark:bg-gray-900 transition-colors duration-200">
		<nav class="bg-white dark:bg-gray-800 shadow-lg">
			<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
<div class="flex justify-between h-16">
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<svg class="h-8 w-auto" viewBox="0 0 200 50" fill="none" xmlns="http://www.w3.org/2000/svg">
								<path d="M10 25C10 16.7157 16.7157 10 25 10H175C183.284 10 190 16.7157 190 25C190 33.2843 183.284 40 175 40H25C16.7157 40 10 33.2843 10 25Z" fill="#4F46E5"/>
								<text x="20" y="32" font-family="Arial, sans-serif" font-size="18" font-weight="bold" fill="white">Codema Patterns</text>
							</svg>
						</div>
						<div class="hidden md:block">
							<div class="ml-10 flex items-baseline space-x-4">
								<a href="/" class="text-gray-800 dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700 px-3 py-2 rounded-md text-sm font-medium">Home</a>
								<a href="/bounties" class="text-gray-800 dark:text-white hover:bg-gray-200 dark:hover:bg-gray-700 px-3 py-2 rounded-md text-sm font-medium">Bounties</a>
							</div>
						</div>
					</div>
					<div class="flex items-center">
						<div class="flex-shrink-0">
							<input type="text" placeholder="Search patterns" class="border rounded-md p-2 dark:bg-gray-700 dark:text-white">
						</div>
						<button id="darkModeToggle" class="ml-4 p-2 rounded-md bg-gray-200 dark:bg-gray-700">
							<svg class="w-6 h-6 text-gray-800 dark:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
							</svg>
						</button>
					</div>
				</div>
			</div>
		</nav>
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
			<h2 class="text-3xl font-extrabold text-gray-900 dark:text-white mb-4">Bounties</h2>
			<p class="text-xl text-gray-600 dark:text-gray-300 mb-8">Explore and contribute to open bounties.</p>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				<!-- Placeholder bounty cards -->
				<div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
					<h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Implement OAuth2 Pattern</h3>
					<p class="text-gray-600 dark:text-gray-400 mb-4">Create a reusable OAuth2 authentication pattern for web applications.</p>
					<span class="inline-block bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-100 text-xs px-2 py-1 rounded-full">$500 Reward</span>
				</div>
				<div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
					<h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Develop GraphQL API Wrapper</h3>
					<p class="text-gray-600 dark:text-gray-400 mb-4">Build a pattern for wrapping RESTful APIs with GraphQL interfaces.</p>
					<span class="inline-block bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-100 text-xs px-2 py-1 rounded-full">$750 Reward</span>
				</div>
				<div class="bg-white dark:bg-gray-800 shadow rounded-lg p-6">
					<h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Create Microservices Boilerplate</h3>
					<p class="text-gray-600 dark:text-gray-400 mb-4">Develop a comprehensive boilerplate for microservices architecture.</p>
					<span class="inline-block bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-100 text-xs px-2 py-1 rounded-full">$1000 Reward</span>
				</div>
			</div>
		</div>
		<script>
			const darkModeToggle = document.getElementById('darkModeToggle');
			darkModeToggle.addEventListener('click', () => {
				document.documentElement.classList.toggle('dark');
				localStorage.theme = document.documentElement.classList.contains('dark') ? 'dark' : 'light';
			});
		</script>
	</body>
	</html>
	`
