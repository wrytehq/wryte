// Form validation utilities
window.FormValidation = {
  isValidEmail: function (email) {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  },

  showError: function (element, message) {
    if (!element) return;
    element.textContent = message;
    element.classList.remove('hidden');
  },

  hideError: function (element) {
    if (!element) return;
    element.textContent = '';
    element.classList.add('hidden');
  },

  clearBackendError: function (inputElement, backendErrorElement) {
    if (backendErrorElement) {
      backendErrorElement.remove();
      inputElement.classList.remove('input-error');
    }
  },

  validateEmail: function (emailInput, emailError) {
    const email = emailInput.value.trim();
    if (email === '') {
      this.hideError(emailError);
      emailInput.classList.remove('input-error');
      return false;
    }
    if (!this.isValidEmail(email)) {
      this.showError(emailError, 'Must be a valid email address');
      emailInput.classList.add('input-error');
      return false;
    }
    this.hideError(emailError);
    emailInput.classList.remove('input-error');
    return true;
  },

  validatePassword: function (passwordInput, passwordError, minLength) {
    const password = passwordInput.value;
    const min = minLength || 6;

    if (password === '') {
      this.hideError(passwordError);
      passwordInput.classList.remove('input-error');
      return false;
    }
    if (password.length < min) {
      this.showError(
        passwordError,
        `Password must be at least ${min} characters`
      );
      passwordInput.classList.add('input-error');
      return false;
    }
    this.hideError(passwordError);
    passwordInput.classList.remove('input-error');
    return true;
  },

  validateConfirmPassword: function (
    passwordInput,
    confirmPasswordInput,
    confirmPasswordError
  ) {
    const password = passwordInput.value;
    const confirmPassword = confirmPasswordInput.value;

    if (confirmPassword === '') {
      this.hideError(confirmPasswordError);
      confirmPasswordInput.classList.remove('input-error');
      return false;
    }
    if (password !== confirmPassword) {
      this.showError(confirmPasswordError, 'Passwords do not match');
      confirmPasswordInput.classList.add('input-error');
      return false;
    }
    this.hideError(confirmPasswordError);
    confirmPasswordInput.classList.remove('input-error');
    return true;
  },

  updateSubmitButton: function (submitBtn, isValid) {
    submitBtn.disabled = !isValid;
    if (isValid) {
      submitBtn.classList.remove('opacity-50', 'cursor-not-allowed');
    } else {
      submitBtn.classList.add('opacity-50', 'cursor-not-allowed');
    }
  },

  setupPasswordToggle: function (
    toggleBtn,
    passwordInput,
    eyeIcon,
    eyeOffIcon
  ) {
    if (!toggleBtn) return;

    toggleBtn.addEventListener('click', function () {
      const type =
        passwordInput.getAttribute('type') === 'password' ? 'text' : 'password';
      passwordInput.setAttribute('type', type);

      eyeIcon.classList.toggle('hidden');
      eyeOffIcon.classList.toggle('hidden');
    });
  },
};
