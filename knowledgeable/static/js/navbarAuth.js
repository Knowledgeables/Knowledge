let cachedUser = undefined;

export async function getCurrentUser() {
    if (cachedUser !== undefined) return cachedUser;

    const res = await fetch("/api/me", {
        credentials: "include"
    });

    cachedUser = res.ok ? await res.json() : null;
    return cachedUser;
}

function baseNav(content) {
    return `
    <div class="flex flex-row justify-between items-center w-full">
     
    <div class="flex items-center gap-14">
              <h1 class="text-[1.5rem] font-semibold tracking-tight">
                Knowledgeable
            </h1>

            <ul class="flex gap-6">
                <li class="rounded-lg hover:bg-white/5 transition">
                    <a href="/" id="nav-home" class="block p-2 text-gray-400 hover:text-gray-300 transition">
                        Home
                    </a>
                </li>

                <li class="rounded-lg hover:bg-white/5 transition">
                    <a href="/search" id="nav-search" class="block p-2 text-gray-400 hover:text-gray-300 transition">
                        Search
                    </a>
                </li>

                <li class="rounded-lg hover:bg-white/5 transition">
                    <a href="#about" id="nav-about" class="block p-2 text-gray-400 hover:text-gray-300 transition">
                        About
                    </a>
                </li>
            </ul>
   </div>


        <div class="flex items-center">
        ${content}
        </div>
    </div>
    `;
}

function guestNav() {
    return `
      <div class="w-64 flex gap-3">

            <a href="/login" class="h-9 w-full px-4 flex items-center justify-center gap-2 rounded-lg
                             text-sm text-gray-400 hover:text-white hover:bg-white/5 transition" id="login-button">
                <i data-lucide="log-in" class="w-4 h-4"></i>
                Login
            </a>

            <a href="/register" class="h-9 w-full px-4 flex items-center justify-center gap-2 rounded-lg
               text-sm font-semibold bg-[#1eff8c] hover:bg-[#3dffa0] 
               text-[#0a0e17] transition" id="nav-register">
                <i data-lucide="user-plus" class="w-4 h-4"></i>
                Sign up
            </a>

        </div>
    `;
}

function userNav(user) {
    return `
    <div class="flex items-center gap-5">
<div class="flex flex-row items-center justify-center gap-2 bg-[#161a20] h-10 px-3 rounded-lg">
    <div   class="h-7 w-7 flex items-center rounded-full justify-center
text-[#ff6b6b]
bg-[#1eff8c]/10
border border-[#1eff8c]/30


transition">
<p>
          ${user.username.substring(0, 1)}
</p>
        </div>

          <p class="text-sm">
            ${user.username}
        </p>
        </div>

<form action="/logout" method="POST">
    <button type="submit"
        class="h-9 w-9 flex items-center justify-center rounded-md
        text-[#ff6b6b]
        bg-[#ff4d4f]/10
        hover:bg-[#ff4d4f]/20
        transition">
        <i data-lucide="log-out" class="w-4 h-4"></i>
    </button>
</form>

    </div>
    `;
}


export async function renderNavbar() {
    const nav = document.getElementById("navbar");
    if (!nav) return;

    // Render a usable guest navbar immediately to avoid flaky selector timing.
    nav.innerHTML = baseNav(guestNav());

    const user = await getCurrentUser();
    if (user) {
        nav.innerHTML = baseNav(userNav(user));
    }

    if (window.lucide) {
        lucide.createIcons();
    }
}

renderNavbar();

