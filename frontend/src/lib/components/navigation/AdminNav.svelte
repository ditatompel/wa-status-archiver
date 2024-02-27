<script>
	import { applyAction, enhance } from '$app/forms';
</script>

<nav class="bg-gray-800 fixed top-0 z-30 w-full shadow-2xl">
	<div class="px-3 py-3 lg:px-5 lg:pl-3">
		<div class="flex items-center justify-between">
			<div class="flex items-center justify-start rtl:justify-end">
				<button
					class="btn btn-sm inline-flex items-center md:hidden"
					aria-label="Mobile Drawer Button"
					on:click={undefined}
				>
					<span>
						<svg viewBox="0 0 100 80" class="fill-token h-4 w-4">
							<rect width="100" height="20" />
							<rect y="30" width="100" height="20" />
							<rect y="60" width="100" height="20" />
						</svg>
					</span>
				</button>

				<a href="/dashboard.html" class="ms-2 flex md:me-24" aria-label="ditatompel">
					<span
						class="text-white hidden self-center whitespace-nowrap text-2xl font-semibold lg:block"
						>ditatompel</span
					>
				</a>
			</div>
			<div class="flex items-center">
				<div class="ms-3 flex items-center space-x-4">
					<form
						action="/logout"
						method="POST"
						use:enhance={async () => {
							return async ({ result }) => {
								await applyAction(result);
							};
						}}
					>
						<input type="hidden" name="logout" value="logout" />
						<button
							class="block w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white"
							role="menuitem"
						>
							Sign out
						</button>
					</form>
				</div>
			</div>
		</div>
	</div>
</nav>
