<script>
  import {
    Navbar,
    NavBrand,
    NavLi,
    NavUl,
    NavHamburger,
    DarkMode,
    Span,
    Button,
    UserCircle,
  } from "flowbite-svelte";
  let SignedIn = false;
  let btnClass =
    "text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg text-sm";
  import { userName, SignIn } from "../stores";

  let user = "";
  userName.subscribe((value) => {
    user = value;
  });
  SignIn.subscribe((value) => {
    SignedIn = value;
  });
</script>

<div>
  <header>
    <Navbar let:hidden let:toggle>
      <NavBrand herf="/">
        <span
          class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
        >
          📝 Notes App
        </span>
      </NavBrand>

      <NavHamburger on:click={toggle} class="dark:text-white" />

      <NavUl {hidden}>
        <NavLi href="/" active={true}>
          Home
        </NavLi>
        {#if SignedIn}
          <NavLi
            href="/history"
            class="self-center whitespace-nowrap dark:text-white"  
          >
          history
          </NavLi>
        {/if}
        <NavLi
          href="/help"
          class="self-center whitespace-nowrap dark:text-white"
        >
          Help
        </NavLi>
        <NavLi href="/about" class="dark:text-white">About</NavLi>
        <NavLi class="dark:text-white">
          <p>👋 {user}</p>
        </NavLi>
        <NavLi href="">
          <DarkMode {btnClass}>
            <svelte:fragment slot="lightIcon">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="1em"
                height="1em"
                preserveAspectRatio="xMidYMid meet"
                viewBox="0 0 32 32"
                ><path
                  fill="currentColor"
                  d="M16 12.005a4 4 0 1 1-4 4a4.005 4.005 0 0 1 4-4m0-2a6 6 0 1 0 6 6a6 6 0 0 0-6-6ZM5.394 6.813L6.81 5.399l3.505 3.506L8.9 10.319zM2 15.005h5v2H2zm3.394 10.193L8.9 21.692l1.414 1.414l-3.505 3.506zM15 25.005h2v5h-2zm6.687-1.9l1.414-1.414l3.506 3.506l-1.414 1.414zm3.313-8.1h5v2h-5zm-3.313-6.101l3.506-3.506l1.414 1.414l-3.506 3.506zM15 2.005h2v5h-2z"
                /></svg
              >
            </svelte:fragment>
            <svelte:fragment slot="darkIcon">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="1em"
                height="1em"
                preserveAspectRatio="xMidYMid meet"
                viewBox="0 0 32 32"
                ><path
                  fill="currentColor"
                  d="M13.502 5.414a15.075 15.075 0 0 0 11.594 18.194a11.113 11.113 0 0 1-7.975 3.39c-.138 0-.278.005-.418 0a11.094 11.094 0 0 1-3.2-21.584M14.98 3a1.002 1.002 0 0 0-.175.016a13.096 13.096 0 0 0 1.825 25.981c.164.006.328 0 .49 0a13.072 13.072 0 0 0 10.703-5.555a1.01 1.01 0 0 0-.783-1.565A13.08 13.08 0 0 1 15.89 4.38A1.015 1.015 0 0 0 14.98 3Z"
                /></svg
              >
            </svelte:fragment>
          </DarkMode>
        </NavLi>
      </NavUl>
    </Navbar>
  </header>

  <main class="">
    <slot />
  </main>
</div>
