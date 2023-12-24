import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import resources from '../resources.vue'

describe('resources', () => {
  it('renders properly', async () => {
    const wrapper = mount(resources, {})
    expect(wrapper.text()).toContain('Resources')
  })
})
