<script>
	import { page } from '$app/stores';
</script>

<aside
	id="logo-sidebar"
	class="bg-gray-800 fixed left-0 top-0 z-20 h-screen w-64 -translate-x-full pt-20 shadow-2xl transition-transform sm:translate-x-0"
	aria-label="Sidebar"
>
	<div class="bg-gray-800 h-full overflow-y-auto px-3 pb-4">
		<ul class="space-y-2 font-medium list-none" data-sveltekit-preload-data="false">
			<li>
				<a
					href="/dashboard.html"
					class={$page.url.pathname === '/dashboard.html' ? 'active' : 'nav-link'}
				>
					<span class="ms-3">Dashboard</span>
				</a>
			</li>
			<li>
				<a
					href="/contacts.html"
					class={$page.url.pathname === '/contacts.html' ? 'active' : 'nav-link'}
				>
					<span class="ms-3">Contacts</span>
				</a>
			</li>
			<li>
				<a
					href="/chats.html"
					class={$page.url.pathname === '/chats.html' ? 'active' : 'nav-link'}
				>
					<span class="ms-3">Chats</span>
				</a>
			</li>
			<li>
				<a
					href="/status-updates.html"
					class={$page.url.pathname === '/status-updates.html' ? 'active' : 'nav-link'}
				>
					<span class="ms-3">Status Updates</span>
				</a>
			</li>
		</ul>
	</div>
</aside>

<style lang="postcss">
	.active {
		@apply flex items-center rounded-lg bg-blue-500 p-2;
	}
	.nav-link {
		@apply text-white flex items-center rounded-lg p-2 hover:bg-pink-500 hover:text-black;
	}
</style>
