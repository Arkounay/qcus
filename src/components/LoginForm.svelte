<script>
    import {getPublicConfig, validatePassword} from '../lib/api.js';
    import {Alert} from "$lib/components/ui/alert";
    import {Input} from "$lib/components/ui/input";
    import * as Field from "$lib/components/ui/field/index.js";
    import {Button} from "$lib/components/ui/button/index.ts";

    let { onlogin } = $props();

    let password = $state('');
    let error = $state('');
    let isLoading = $state(false);
    let showDefaultPasswordHint = $state(true);
    let passwordInput = $state(null);

    // Fetch config once (cached at module level in api.js)
    getPublicConfig().then(config => {
        showDefaultPasswordHint = config.isDefaultPassword;
    });

    $effect(() => {
        if (error && passwordInput) {
            passwordInput.focus();
        }
    });

    async function handleSubmit() {
        if (!password) {
            error = 'Please enter a password.';
            return;
        }

        isLoading = true;
        error = '';

        const isValid = await validatePassword(password);

        if (isValid) {
            onlogin?.({password});
            password = '';
        } else {
            error = 'Invalid password. Please try again.';
            password = '';
        }

        isLoading = false;
    }
</script>

<div class="flex min-h-screen items-center justify-center p-4 bg-gray-100">
    <div class="w-full max-w-md rounded-lg border border-border bg-card shadow-lg bg-white">
        <div class="p-6">
            <h2 class="text-2xl font-semibold text-card-foreground">File Upload Server</h2>
            <p class="mt-2 text-sm text-muted-foreground">
                Please enter the upload password to continue </p>

            <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="mt-6 space-y-4">
                <!--<div>
                    <input type="password" bind:value={password} placeholder="Enter password" disabled={isLoading} class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50" class:border-destructive={error}/>
                </div>-->
                <Field.Set>
                    <Field.Group>
                        <Field.Field data-invalid={error !== ''}>
                            <Input bind:ref={passwordInput} type="password" bind:value={password} placeholder="Enter password" disabled={isLoading} aria-invalid={error !== ''}/>
                            {#if error}
                                <Field.Error> {error}</Field.Error>
                            {/if}
                        </Field.Field>
                    </Field.Group>
                </Field.Set>

                <div class="flex justify-end">
                    <Button type="submit" disabled={isLoading} class="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50">
                        {#if isLoading}
                            <svg class="mr-2 h-4 w-4 animate-spin" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            Validating...
                        {:else}
                            Login
                        {/if}
                    </Button>
                </div>
            </form>
        </div>

        {#if showDefaultPasswordHint}
            <div class="border-t border-border bg-muted/50 p-4">
                <div class="flex gap-2 text-sm">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="h-5 w-5 shrink-0 text-muted-foreground">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z"/>
                    </svg>
                    <div class="text-muted-foreground">
                        <strong>Default password:</strong> <pre class="inline-block ml-1 rounded bg-background px-1.5 py-0.5 font-mono text-xs font-semibold">demo</pre>
                    </div>
                </div>
            </div>
        {/if}
    </div>
</div>
