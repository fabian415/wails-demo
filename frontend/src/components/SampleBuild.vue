<template>
  <div class="sample-build-test">
    <div class="container">
      <SampleBuildProgress />
    </div>
  </div>
</template>

<script>
import SampleBuildProgress from '@/components/SampleBuildProgress.vue'

export default {
  name: 'SampleBuildTest',
  components: {
    SampleBuildProgress
  }
}
</script>

<style scoped>
.sample-build-test {
  min-height: 100vh;
  /* background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px; */
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

h1 {
  text-align: center;
  color: #333;
  margin: 0;
  padding: 30px 20px 10px;
  font-size: 2.5rem;
  font-weight: 300;
}

.description {
  text-align: center;
  color: #666;
  margin: 0 0 30px;
  padding: 0 20px;
  font-size: 1.1rem;
  line-height: 1.6;
}

.info-section {
  background: #f8f9fa;
  padding: 30px;
  margin-top: 20px;
  border-top: 1px solid #e9ecef;
}

.info-section h3 {
  color: #495057;
  margin: 0 0 15px 0;
  font-size: 1.3rem;
  font-weight: 600;
}

.info-section ul {
  margin: 0 0 25px 0;
  padding-left: 20px;
}

.info-section li {
  margin-bottom: 8px;
  color: #6c757d;
  line-height: 1.5;
}

.info-section code {
  background: #e9ecef;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9em;
  color: #e83e8c;
}

strong {
  color: #495057;
  font-weight: 600;
}

/* Responsive design */
@media (max-width: 768px) {
  .sample-build-test {
    padding: 10px;
  }
  
  h1 {
    font-size: 2rem;
    padding: 20px 15px 10px;
  }
  
  .description {
    padding: 0 15px;
    font-size: 1rem;
  }
  
  .info-section {
    padding: 20px 15px;
  }
}
</style>
